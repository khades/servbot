package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

func GetTwitchDJEnabledChannels() (*[]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := Db.C(channelInfoCollection).Find(bson.M{"twitchdj.id": bson.M{"$exists": true, "$ne": ""}}).All(&result)
	return &result, error
}
