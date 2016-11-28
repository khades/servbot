package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

func GetUserMessageHistory(user *string, channel *string) (*models.ChatMessageLog, error) {
	result := models.ChatMessageLog{}
	error := Db.C("messageLogs").Find(bson.M{"channel": *channel, "user": *user}).One(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}
