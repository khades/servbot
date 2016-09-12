package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// PushBanme sets if banme command enabled
func PushBanme(channel string, banme models.Banme) {
	channelInfo, _ := GetChannelInfo(channel)
	if channelInfo != nil {
		channelInfo.Banme = banme
	} else {
		channelInfoRepositoryObject.forceCreateObject(channel, &models.ChannelInfo{Channel: channel, Banme: banme})
	}
	Db.C("channelInfo").Upsert(models.ChannelSelector{Channel: channel}, bson.M{"$set": bson.M{"banme": banme}})
}
