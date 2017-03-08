package repos

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushStreamStatus updates stream status (start of stream, topic of stream)
func PushStreamStatus(channelID *string, streamStatus *models.StreamStatus) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.StreamStatus = *streamStatus
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, StreamStatus: *streamStatus})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"streamstatus": *streamStatus}})
}
