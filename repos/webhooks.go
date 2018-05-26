package repos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/khades/servbot/utils"
	"github.com/sirupsen/logrus"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
)

var webhooklibrary = "webhooklibrary"

// func PutChallengeForWebHookTopic(channelID *string, topic *string, challenge *string) {
// 	db.C(channelInfoCollection).Update(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"challenge": *challenge}})
// }

type hub struct {
	Mode         string `json:"hub.mode"`
	Topic        string `json:"hub.topic"`
	Callback     string `json:"hub.callback"`
	LeaseSeconds string `json:"hub.lease_seconds"`
	Secret       string `json:"hub.secret"`
}

func updateWebHookTopic(channelID *string, topic string, secret *string, expiresAt time.Time) {
	db.C(webhooklibrary).Upsert(bson.M{"channelid": *channelID, "topic": topic}, bson.M{"$set": bson.M{"secret": *secret, "expiresat": expiresAt}})
}

func GetWebHookTopic(channelID *string, topic string) (*models.WebHookInfo, error) {
	var result models.WebHookInfo
	err := db.C(webhooklibrary).Find(bson.M{"channelid": *channelID, "topic": topic}).One(&result)
	return &result, err
}

func getHooksForChannel(channelID *string) ([]models.WebHookInfo, error) {
	var result []models.WebHookInfo
	err := db.C(webhooklibrary).Find(bson.M{"channelid": *channelID}).All(&result)
	return result, err
}

func getNonExpiredHooks(pollDuration time.Duration) ([]models.WebHookInfo, error) {
	var result []models.WebHookInfo
	err := db.C(webhooklibrary).Find(bson.M{"expiresat": bson.M{"$gte": time.Now().Add(pollDuration)}}).All(&result)
	return result, err
}

func CheckAndSubscribeToWebhooks(pollDuration time.Duration) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "webhook",
		"action":  "CheckAndSubscribeToWebhooks"})
	logger.Debugf("Starting")

	channels, error := GetActiveChannels()
	if error != nil {
		logger.Debugf("Error %s", error.Error())

		return
	}
	logger.Debugf("Channels found: %d", len(channels))

	nonExpiredHooks, _ := getNonExpiredHooks(pollDuration)
	for _, channel := range channels {
		logger.Debugf("Processing channel %s", channel.ChannelID)

		followsFound, _ := getExpiredTopics(nonExpiredHooks, channel.ChannelID)
		if followsFound == false {
			SubChannelToFollowerHooks(channel.ChannelID)
		}
	}
}

func SubChannelToFollowerHooks(channelID string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "webhook",
		"action":  "SubChannelToFollowerHooks"})
	secret := utils.RandomString(10)
	form := hub{
		Mode:         "subscribe",
		Topic:        "https://api.twitch.tv/helix/users/follows?to_id=" + channelID,
		Callback:     "https://servbot.khades.org/api/webhook/follows",
		LeaseSeconds: "864000",
		Secret:       secret}

	body, _ := json.Marshal(form)

	resp, _ := twitchHelixPost("webhooks/hub", bytes.NewReader(body))
	if resp != nil {
		defer resp.Body.Close()
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err == nil {
		logger.Debugf("Repsonse is %q", dump)
	}
	logger.Debugf("Status is %d", resp.StatusCode)
	if resp.StatusCode == http.StatusAccepted {
		logger.Debugf("Doing update for channel %s", channelID)
		updateWebHookTopic(&channelID, "follows", &secret, time.Now().Add(10*24*time.Hour))
	}
}

func getExpiredTopics(nonExpiredTopics []models.WebHookInfo, channelID string) (bool, bool) {
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
