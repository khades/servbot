package repos

import (
	"time"

	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var messageLogsCollection = "messageLogs"

// GetChannelUsers returns list of all users, who even wrote in chat room, with optional pattern to match
func GetChannelUsers(channelID *string, pattern *string) ([]models.ChannelUser, error) {
	var channelUsers []models.ChannelUser
	if *pattern == "" {
		error := db.C(messageLogsCollection).Find(models.ChannelSelector{ChannelID: *channelID}).Sort("messages.date").Limit(100).All(&channelUsers)
		return channelUsers, error

	}
	error := db.C(messageLogsCollection).Find(bson.M{
		"channelid": *channelID,
		"knownnicknames": bson.M{
			"$regex":   *pattern,
			"$options": "i"}}).Sort("messages.date").Limit(100).All(&channelUsers)
	return channelUsers, error
}

//LogMessage logs the users chat message on channel with logging its known nicknames
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
	db.C(messageLogsCollection).Upsert(
		bson.M{"channelid": message.ChannelID, "userid": message.UserID},
		query)
}

// GetUserMessageHistoryByUserID returns message history for user, specified by userID on specific channel
func GetUserMessageHistoryByUserID(userID *string, channelID *string) (*models.ChatMessageLog, error) {
	result := models.ChatMessageLog{}
	error := db.C("messageLogs").Find(bson.M{"channelid": *channelID, "userid": *userID}).One(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}

// GetUserMessageHistoryByUsername returns message history for user, specified by displayname on specific channel
func GetUserMessageHistoryByUsername(user *string, channelID *string) (*models.ChatMessageLog, error) {
	result := models.ChatMessageLog{}
	error := db.C("messageLogs").Find(bson.M{"channelid": *channelID, "user": *user}).One(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}

// GetUserMessageHistoryByKnownUsernames returns message history for user, specified by his knows display names on specific channel
func GetUserMessageHistoryByKnownUsernames(user *string, channelID *string) ([]models.ChatMessageLog, error) {
	result := []models.ChatMessageLog{}
	error := db.C("messageLogs").Find(bson.M{"channelid": *channelID, "knownusernames": *user}).All(&result)
	if error != nil {
		return result, error
	}
	return result, error
}
