package repos

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushSubAlert updates stream status (start of stream, topic of stream)
func PushSubAlert(channel string, subAlert models.SubAlert) {
	channelInfo, _ := GetChannelInfo(channel)
	if channelInfo != nil {
		channelInfo.SubAlert = subAlert
	} else {
		channelInfoRepositoryObject.forceCreateObject(channel, &models.ChannelInfo{Channel: channel, SubAlert: subAlert})
	}
	Db.C("channelInfo").Upsert(models.ChannelSelector{Channel: channel}, bson.M{"$set": bson.M{"subAlert": subAlert}})
}
