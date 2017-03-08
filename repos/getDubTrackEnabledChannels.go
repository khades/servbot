package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

func GetDubTrackEnabledChannels() (*[]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := Db.C(channelInfoCollection).Find(bson.M{"dubtrack.id": bson.M{"$exists": true, "$ne": ""}}).All(&result)
	return &result, error
}
