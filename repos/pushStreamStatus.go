package repos

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushStreamStatus updates list of mods on channel
func PushStreamStatus(channel string, streamStatus models.StreamStatus) {
	channelInfo, _ := GetChannelInfo(channel)
	if channelInfo != nil {
		channelInfo.StreamStatus = &streamStatus
	} else {
		channelInfoRepositoryObject.forceCreateObject(channel, &models.ChannelInfo{Channel: channel, StreamStatus: &streamStatus})
	}
	Db.C("channelInfo").Upsert(models.ChannelSelector{Channel: channel}, bson.M{"$set": bson.M{"streamStatus": streamStatus}})
}
