package repos

// import (
// 	"github.com/khades/servbot/models"
// 	"gopkg.in/mgo.v2/bson"
// )

// // PushBanme sets if banme command enabled
// func PushBanme(channelID *string, banme *models.Banme) {
// 	channelInfo, _ := GetChannelInfo(channelID)
// 	if channelInfo != nil {
// 		channelInfo.Banme = *banme
// 	} else {
// 		channelInfoRepositoryObject.forceCreateObject(*channel, &models.ChannelInfo{ChannelID: *channelID, Banme: *banme})
// 	}
// 	Db.C("channelInfo").Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"banme": *banme}})
// }
