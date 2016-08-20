package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

//LogMessage logs the chat message for showing it later
func LogMessage(message models.ChatMessage) {
	Db.C("messageLogs").Insert(message)
	Db.C("channelUsers").Upsert(models.ChannelSelector{Channel: message.Channel}, bson.M{"$addToSet": bson.M{"users": message.User}})
}
