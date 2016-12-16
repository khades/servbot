package repos

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushDubTrack updates stream DubTrack status (is it playing or is there any songs)
func PushDubTrack(channel *string, dubTrack *models.DubTrack) {
	channelInfo, _ := GetChannelInfo(channel)
	if channelInfo != nil {
		channelInfo.DubTrack = *dubTrack
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channel, &models.ChannelInfo{Channel: *channel, DubTrack: *dubTrack})
	}
	Db.C("channelInfo").Upsert(models.ChannelSelector{Channel: *channel}, bson.M{"$set": bson.M{"dubtrack": *dubTrack}})
}
