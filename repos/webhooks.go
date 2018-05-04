package repos

import (
	"time"

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

func upsateWebHookTopic(channelID *string, topic string, secret *string, expiresAt time.Time) {
	db.C(channelInfoCollection).Upsert(bson.M{"channelid": *channelID, "topic": topic}, bson.M{"$set": bson.M{"secret": *secret, "expiresat": expiresAt}})
}

func GetWebHookTopic(channelID *string, topic string) (*models.WebHookInfo, error) {
	var result models.WebHookInfo
	err := db.C(channelInfoCollection).Find(bson.M{"channelid": *channelID, "topic": topic}).One(&result)
	return &result, err
}

func getHooksForChannel(channelID *string) ([]models.WebHookInfo, error) {
	var result []models.WebHookInfo
	err := db.C(channelInfoCollection).Find(bson.M{"channelid": *channelID}).All(&result)
	return result, err
}

func getNonExpiredHooks(pollDuration time.Duration) ([]models.WebHookInfo, error) {
	var result []models.WebHookInfo
	err := db.C(channelInfoCollection).Find(bson.M{"expiresat": bson.M{"$lte": time.Now().Add(-pollDuration)}}).All(&result)
	return result, err
}

// func CheckAndSubscribeToWebhooks(pollDuration time.Duration) {
// 	channels, error := GetActiveChannels()
// 	if error != nil {
// 		return
// 	}
// 	nonExpiredHooks, _ := getNonExpiredHooks(pollDuration)
// 	for _, channel := range channels {
// 		followsFound, streamsFound := getExpiredTopics(nonExpiredHooks, channel.ChannelID)
// 		secret:= "s3cR37"
// 		form := hub{
// 			Mode:"subscibe",
// 			Topic:"https://api.twitch.tv/helix/users/follows?to_id="+channel.ChannelID,
// 			Callback: "https://servbot.khades.org/api/webhook/follows",
// 			LeaseSeconds: "864000",
// 			Secret: secret		}

// 		upsateWebHookTopic(&channel.ChannelID, "follows", &secret, time.Now().Add(10*24*time.Hour))
// 		twitchHelixPost("webhooks/hub", form.Encode())
// 	}
// }

func getExpiredTopics(nonExpiredTopics []models.WebHookInfo, channelID string) (bool, bool) {
	if len(nonExpiredTopics) == 0 {
		return  false, false
	}
	followsFound := false
	streamsFound := false
	for _, topic := range nonExpiredTopics {
		if (topic.ChannelID == channelID) {
			if (topic.Topic == "follows") {
				followsFound= true
			}
			if (topic.Topic == "streams") {
				streamsFound= true
			}
			if (followsFound == true && streamsFound == true){
				break
			}
		}
	}
	return followsFound, streamsFound
}