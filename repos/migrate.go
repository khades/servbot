package repos

import "gopkg.in/mgo.v2/bson"

func Migrate() {
	channelsWithID, error := GetUsersID(&Config.Channels)
	if error != nil {
		return
	}
	Db.C("subAlertHistory").DropCollection()
	Db.C("messageLogs").DropCollection()
	for channel, channelID := range *channelsWithID {
		Db.C(subAlertCollection).UpdateAll(bson.M{"channel": channel}, bson.M{
			"$set":   bson.M{"channelid": channelID},
			"$unset": bson.M{"history": ""}})
		Db.C(channelInfoCollection).UpdateAll(bson.M{"channel": channel}, bson.M{
			"$set":   bson.M{"channelid": channelID},
			"$unset": bson.M{"history": ""}})
		Db.C(templateCollection).UpdateAll(bson.M{"channel": channel}, bson.M{
			"$set":   bson.M{"channelid": channelID},
			"$unset": bson.M{"history": ""}})
	}
}
