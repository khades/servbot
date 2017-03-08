package repos

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushTwitchDJ updates stream twitchDJ status (is it playing or is there any songs)
func PushTwitchDJ(channelID *string, twitchDJ *models.TwitchDJ) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.TwitchDJ = *twitchDJ
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, TwitchDJ: *twitchDJ})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"twitchdj": *twitchDJ}})
}
