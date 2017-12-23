package repos

import (
	"time"

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
			"$slice": 50}}}
	if message.MessageType == "timeout" || message.MessageType == "ban" {
		banInfo := models.BanInfo{User: message.User,
			Duration: message.BanLength,
			Type:     message.MessageType,
			Date:     time.Now()}
		query = bson.M{
			"$set":      bson.M{"user": message.User, "channel": message.Channel},
			"$addToSet": bson.M{"knownnicknames": message.User},
			"$push": bson.M{"messages": bson.M{
				"$each":  []models.MessageStruct{message.MessageStruct},
				"$sort":  bson.M{"date": -1},
				"$slice": 50},
				"bans": bson.M{
					"$each":  []models.BanInfo{banInfo},
					"$sort":  bson.M{"date": -1},
					"$slice": 10}}}
		LogChannelBan(&message.UserID, &message.User, &message.ChannelID, &message.BanLength)
	}
	Db.C("messageLogs").Upsert(
		bson.M{"channelid": message.ChannelID, "userid": message.UserID},
		query)
}
