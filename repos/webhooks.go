package repos

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
)

var webhooklibrary = "webhooklibrary"

func PutChallengeForWebHookTopic(channelID *string, topic *string, challenge *string) {
	db.C(channelInfoCollection).Update(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"challenge": *challenge}})
}

func PutSecretForWebHookTopic(channelID *string, topic *string, secret *string) {
	db.C(channelInfoCollection).Update(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"secret": *secret}})
}

func putTimeoutForWebHookTopic(channelID *string, topic *string, expiresAt time.Time) {
	db.C(channelInfoCollection).Update(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"expiresat": expiresAt}})
}

func GetWebHookTopic(channelID *string, topic *string) (*models.WebHookInfo, error) {
	var result models.WebHookInfo
	err := db.C(channelInfoCollection).Find(bson.M{"channelid": *channelID, "topic": *topic}).One(&result)
	return &result, err
}

func getHooksForChannel(channelID *string) (*models.WebHookInfo, error) {
	var result models.WebHookInfo
	err := db.C(channelInfoCollection).Find(bson.M{"channelid": *channelID}).One(&result)
	return &result, err
}

// func CheckAndSubscribeToWebhooks() {
// 	channels, error := GetActiveChannels()
// 	if error != nil {
// 		return
// 	}
// 	for channel, _ := range channels {

// 	}
// }
