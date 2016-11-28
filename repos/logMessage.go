package repos

import (
	"log"

	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

//LogMessage logs the chat message
func LogMessage(message *models.ChatMessage) {
	log.Println("Logging message")
	query := bson.M{
		"$push": bson.M{"messages": bson.M{
			"$each":  []models.MessageStruct{message.MessageStruct},
			"$sort":  bson.M{"date": -1},
			"$slice": 25}}}
	Db.C("messageLogs").Upsert(
		bson.M{"channel": message.Channel, "user": message.User},
		query)
}
