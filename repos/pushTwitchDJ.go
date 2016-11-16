package repos

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushTwitchDJ updates stream twitchDJ status (is it playing or is there any songs)
func PushTwitchDJ(channel *string, twitchDJ *models.TwitchDJ) {
	channelInfo, _ := GetChannelInfo(channel)
	if channelInfo != nil {
		channelInfo.TwitchDJ = *twitchDJ
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channel, &models.ChannelInfo{Channel: *channel, TwitchDJ: *twitchDJ})
	}
	Db.C("channelInfo").Upsert(models.ChannelSelector{Channel: *channel}, bson.M{"$set": bson.M{"twitchdj": *twitchDJ}})
}
