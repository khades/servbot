package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// PushMods updates list of mods on channel
func PushMods(channelID *string, mods *[]string) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.Mods = *mods
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, Mods: *mods})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"mods": *mods}})
}
