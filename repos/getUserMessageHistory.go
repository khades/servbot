package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

func GetUserMessageHistoryByUserID(userID *string, channelID *string) (*models.ChatMessageLog, error) {
	result := models.ChatMessageLog{}
	error := Db.C("messageLogs").Find(bson.M{"channelid": *channelID, "userid": *userID}).One(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}

func GetUserMessageHistoryByUsername(user *string, channelID *string) (*models.ChatMessageLog, error) {
	result := models.ChatMessageLog{}
	error := Db.C("messageLogs").Find(bson.M{"channelid": *channelID, "user": *user}).One(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}

func GetUserMessageHistoryByKnownUsernames(user *string, channelID *string) (*[]models.ChatMessageLog, error) {
	result := []models.ChatMessageLog{}
	error := Db.C("messageLogs").Find(bson.M{"channelid": *channelID, "knownusernames": *user}).All(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}
