package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

//LogMessage logs the chat message
func LogMessage(message *models.ChatMessage) {
	query := bson.M{
		"$set":      bson.M{"user": message.User, "channel": message.Channel},
		"$addToSet": bson.M{"knownnicknames": message.User},
		"$push": bson.M{"messages": bson.M{
			"$each":  []models.MessageStruct{message.MessageStruct},
			"$sort":  bson.M{"date": -1},
			"$slice": 20}}}
	Db.C("messageLogs").Upsert(
		bson.M{"channelid": message.ChannelID, "userid": message.UserID},
		query)
}
