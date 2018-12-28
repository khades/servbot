package webhook

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchAPI"
	"github.com/khades/servbot/utils"
	"github.com/sirupsen/logrus"
)

// func PutChallengeForWebHookTopic(channelID *string, topic *string, challenge *string) {
// 	db.C(channelInfoCollection).Update(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"challenge": *challenge}})
// }

type Service struct {
	collection         *mgo.Collection
	channelInfoService *channelInfo.Service
	twitchAPIService   *twitchAPI.Client
}

func (service *Service) update(channelID *string, topic string, secret *string, expiresAt time.Time) {
	service.collection.Upsert(bson.M{"channelid": *channelID, "topic": topic}, bson.M{"$set": bson.M{"secret": *secret, "expiresat": expiresAt}})
}

func (service *Service) Get(channelID *string, topic string) (*WebHookInfo, error) {
	var result WebHookInfo
	err := service.collection.Find(bson.M{"channelid": *channelID, "topic": topic}).One(&result)
	return &result, err
}

func (service *Service) list(channelID *string) ([]WebHookInfo, error) {
	var result []WebHookInfo
	err := service.collection.Find(bson.M{"channelid": *channelID}).All(&result)
	return result, err
}

func (service *Service) getNonExpired(pollDuration time.Duration) ([]WebHookInfo, error) {
	var result []WebHookInfo
	err := service.collection.Find(bson.M{"expiresat": bson.M{"$gte": time.Now().Add(pollDuration)}}).All(&result)
	return result, err
}

func (service *Service) Subscribe(pollDuration time.Duration) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "webhook",
		"action":  "Subscribe"})

	channels, error := service.channelInfoService.GetActiveChannels()
	if error != nil {
		logger.Debugf("Error %s", error.Error())

		return
	}
	logger.Debugf("Channels found: %d", len(channels))

	nonExpiredHooks, _ := service.getNonExpired(pollDuration)
	for _, channel := range channels {
		logger.Debugf("Processing channel %s", channel.ChannelID)

		followsFound, _ := service.getExpired(nonExpiredHooks, channel.ChannelID)
		if followsFound == false {
			service.subToFollowerHook(channel.ChannelID)
		}
	}
}

func (service *Service) subToFollowerHook(channelID string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "webhook",
		"action":  "SubToFollowerHook"})
	secret := utils.RandomString(10)
	success := service.twitchAPIService.SubscribeToChannelFollowerWebhook(channelID, secret)
	if success == true {
		logger.Debugf("Doing update for channel %s", channelID)
		service.update(&channelID, "follows", &secret, time.Now().Add(10*24*time.Hour))
	} else {
		logger.Debugf("Follower hook is not updated for channel %s", channelID)

	}
}

func (service *Service) getExpired(nonExpiredTopics []WebHookInfo, channelID string) (bool, bool) {
	if len(nonExpiredTopics) == 0 {
		return false, false
	}
	followsFound := false
	streamsFound := false
	for _, topic := range nonExpiredTopics {
		if topic.ChannelID == channelID {
			if topic.Topic == "follows" {
				followsFound = true
			}
			if topic.Topic == "streams" {
				streamsFound = true
			}
			if followsFound == true && streamsFound == true {
				break
			}
		}
	}
	return followsFound, streamsFound
}
