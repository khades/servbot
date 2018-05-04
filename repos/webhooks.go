package repos

import (
	"net/url"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
)

var webhooklibrary = "webhooklibrary"

// func PutChallengeForWebHookTopic(channelID *string, topic *string, challenge *string) {
// 	db.C(channelInfoCollection).Update(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"challenge": *challenge}})
// }

type hub struct {
	Mode string `json:"hub.mode"`
	Topic string `json:"hub.topic"`
	Callback string `json:"hub.callback"`
	LeaseSeconds string `json:"hub.lease_seconds"`
	Secret string `json:"hub.secret"`
}

func upsateWebHookTopic(channelID *string, topic *string, secret *string, expiresAt time.Time) {
	db.C(channelInfoCollection).Upsert(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"secret": *secret, "expiresat": expiresAt}})
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

// func CheckAndSubscribeToWebhooks(pollDuration time.Duration) {
// 	channels, error := GetActiveChannels()
// 	if error != nil {
// 		return
// 	}
// 	for _, channel := range channels {
// 		item, itemError := GetWebHookTopic(&channel.ChannelID, "follows")
// 		if itemError == nil && item.ExpiresAt.Sub(time.Now()) > pollDuration {
// 			continue
// 		}
// 		form := hub{
// 			Mode:"subscibe",
// 			Topic:"https://api.twitch.tv/helix/users/follows?to_id="+channel.ChannelID,
// 			Callback: "https://servbot.khades.org/api/webhook/follows",
// 			LeaseSeconds: "864000",
// 			Secret: "s3cR37"		}

// 		twitchHelixPost("webhooks/hub", form.Encode())
// 	}
// }
