package repos

import "gopkg.in/mgo.v2/bson"

func CreateChannels() {
	channelsWithID, error := GetUsersID(&Config.Channels)
	if error != nil {
		return
	}
	for channel, channelID := range *channelsWithID {
		Db.C(channelInfoCollection).Upsert(bson.M{"channelid": channelID}, bson.M{
			"$set": bson.M{"channel": channel}})
	}
}
