package repos

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushDubTrack updates stream DubTrack status (is it playing or is there any songs)
func PushDubTrack(channelID *string, dubTrack *models.DubTrack) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.DubTrack = *dubTrack
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, DubTrack: *dubTrack})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"dubtrack": *dubTrack}})
}
