package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// GetTwitchDJEnabledChannels returns list of all twitchdj-configured channels
func GetTwitchDJEnabledChannels() ([]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := db.C(channelInfoCollection).Find(bson.M{"twitchdj.id": bson.M{"$exists": true, "$ne": ""}}).All(&result)
	return result, error
}

// PushTwitchDJ updates stream twitchDJ status (is it playing or is there any songs)
func PushTwitchDJ(channelID *string, twitchDJ *models.TwitchDJ) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.TwitchDJ = *twitchDJ
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, TwitchDJ: *twitchDJ})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"twitchdj": *twitchDJ}})
}
