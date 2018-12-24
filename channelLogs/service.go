package channelLogs

import (
	"time"

	"github.com/khades/servbot/chatMessage"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelBans"
	"github.com/khades/servbot/userResolve"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
)

type Service struct {
	// Dependencies
	collection         *mgo.Collection
	channelBansService *channelBans.Service
	userResolveService *userResolve.Service
}

// GetUsers returns list of all users, who even wrote in chat room, with optional pattern to match
func (service *Service) GetUsers(channelID *string, pattern *string) ([]models.ChannelUser, error) {
	var channelUsers []models.ChannelUser
	if *pattern == "" {
		error := service.collection.Find(models.ChannelSelector{ChannelID: *channelID}).Sort("-lastupdate").Limit(100).All(&channelUsers)
		return channelUsers, error

	}
	error := service.collection.Find(bson.M{
		"channelid": *channelID,
		"knownnicknames": bson.M{
			"$regex":   *pattern,
			"$options": "i"}}).Sort("-lastupdate").Limit(100).All(&channelUsers)
	return channelUsers, error
}

//Log logs the users chat message on channel with logging its known nicknames
func (service *Service) Log(message *chatMessage.ChatMessage) {
	if message.UserID != "" && message.UserID != "" {
		service.userResolveService.Update(&message.UserID, &message.User)
	}

	query := bson.M{
		"$set":      bson.M{"user": message.User, "channel": message.Channel, "lastupdate": time.Now()},
		"$addToSet": bson.M{"knownnicknames": message.User},
		"$push": bson.M{"messages": bson.M{
			"$each":  []chatMessage.MessageStruct{message.MessageStruct},
			"$sort":  bson.M{"date": -1},
			"$slice": 50}}}
	if message.MessageType == "timeout" || message.MessageType == "ban" || message.MessageType == "unban" || message.MessageType == "untimeout" {
		banInfo := BanInfo{User: message.User,
			Duration:    message.BanLength,
			Type:        message.MessageType,
			Reason:      message.BanReason,
			BanIssuer:   message.BanIssuer,
			BanIssuerID: message.BanIssuerID,
			Date:        time.Now()}
		query = bson.M{
			"$set":      bson.M{"user": message.User, "channel": message.Channel, "lastupdate": time.Now()},
			"$addToSet": bson.M{"knownnicknames": message.User},
			"$push": bson.M{"messages": bson.M{
				"$each":  []chatMessage.MessageStruct{message.MessageStruct},
				"$sort":  bson.M{"date": -1},
				"$slice": 50},
				"bans": bson.M{
					"$each":  []BanInfo{banInfo},
					"$sort":  bson.M{"date": -1},
					"$slice": 10}}}
		if message.MessageType == "timeout" || message.MessageType == "ban" {
			service.channelBansService.Log(&message.UserID, &message.User, &message.ChannelID, &message.BanLength)
		}
	}
	service.collection.Upsert(
		bson.M{"channelid": message.ChannelID, "userid": message.UserID},
		query)
}

// GetByUserID returns message history for user, specified by userID on specific channel
func (service *Service) GetByUserID(userID *string, channelID *string) (*ChatMessageLog, error) {
	result := ChatMessageLog{}
	error := service.collection.Find(bson.M{"channelid": *channelID, "userid": *userID}).One(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}

// GetByUserName returns message history for user, specified by displayname on specific channel
func (service *Service) GetByUserName(user *string, channelID *string) (*ChatMessageLog, error) {
	result := ChatMessageLog{}
	error := service.collection.Find(bson.M{"channelid": *channelID, "user": *user}).One(&result)
	if error != nil {
		return &result, error
	}
	return &result, error
}

// GetByKnownNicknames returns message history for user, specified by his knows display names on specific channel
func (service *Service) GetByKnownNicknames(user *string, channelID *string) ([]ChatMessageLog, error) {
	result := []ChatMessageLog{}
	error := service.collection.Find(bson.M{"channelid": *channelID, "knownusernames": *user}).All(&result)
	if error != nil {
		return result, error
	}
	return result, error
}
