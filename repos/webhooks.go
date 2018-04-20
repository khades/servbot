package repos

import (
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
)

var webhooklibrary = "webhooklibrary"

func PutChallengeForWebHookTopic(channelID *string, topic *string, challenge *string) {
	db.C(channelInfoCollection).Update(bson.M{"channelid": *channelID, "topic": *topic}, bson.M{"$set": bson.M{"challenge": *challenge}})
}

func GetChallengeForWebHookTopic(channelID *string, topic *string) (string, bool) {
	var result models.WebHookInfo
	err := db.C(channelInfoCollection).Find(bson.M{"channelid": *channelID, "topic": *topic}).One(&result)
	if err != nil {
		return "", false
	}
	return result.Challenge, true
}

func CheckAndSubscribeToWebhooks() {

}
