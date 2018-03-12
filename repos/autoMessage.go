package repos

import (
	"errors"
	"strings"
	"time"

	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

var autoMessageCollectionName = "autoMessages"

// DecrementAutoMessages decrement message threshold every time someone writes in chat, with check if automessage is bound to specific game on stream
func DecrementAutoMessages(channelInfo *models.ChannelInfo) {
	games := []string{""}
	if channelInfo.StreamStatus.Online == true {
		games = append(games, channelInfo.StreamStatus.Game)
	}
	db.C(autoMessageCollectionName).UpdateAll(bson.M{
		"channelid": channelInfo.ChannelID,
		"message":   bson.M{"$ne": ""},
		"$or": []bson.M{
			bson.M{"game": bson.M{"$in": games}},
			bson.M{"game": bson.M{"$exists": false}}}},
		bson.M{"$inc": bson.M{"messagethreshold": -1}})
}

// RemoveInactiveAutoMessages removes all automessages on channel, which had no update for a week
func RemoveInactiveAutoMessages(channelID *string) {
	db.C(autoMessageCollectionName).RemoveAll(bson.M{
		"channelid":    *channelID,
		"message":      "",
		"history.date": bson.M{"$not": bson.M{"$gte": time.Now().Add(24 * -7 * time.Hour)}}})
}

// GetCurrentAutoMessages returns list of ALL automessage, which are served on that chatbot instance
func GetCurrentAutoMessages() ([]models.AutoMessage, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "automessages",
		"action":  "GetCurrentAutoMessage"})
	logger.Debug("AutoMessage: Getting Current AutoMessages")
	var result []models.AutoMessage
	error := db.C(autoMessageCollectionName).Find(bson.M{
		"message":           bson.M{"$ne": ""},
		"messagethreshold":  bson.M{"$lte": 0},
		"durationthreshold": bson.M{"$lte": time.Now()}}).All(&result)
	logger.Infof("AutoMessage: Got %d AutoMessages", len(result))
	return result, error
}

// ResetAutoMessageThreshold resets automessage thresholds after successfull automessage execution
func ResetAutoMessageThreshold(autoMessage *models.AutoMessage) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "automessages",
		"action":  "ResetAutoMessageThreshold"})
	logger.Infof("AutoMessage: Resetting AutoMessage %s", autoMessage.ID)
	now := time.Now()
	db.C(autoMessageCollectionName).Update(bson.M{"_id": autoMessage.ID}, bson.M{"$set": bson.M{
		"messagethreshold":  autoMessage.MessageLimit,
		"durationthreshold": now.Add(autoMessage.DurationLimit)}})
}

// GetAutoMessage returns specific automessage for specific channel
func GetAutoMessage(id *string, channelID *string) (*models.AutoMessageWithHistory, error) {
	var result models.AutoMessageWithHistory
	objectID := bson.ObjectIdHex(*id)
	error := db.C(autoMessageCollectionName).Find(bson.M{"_id": objectID, "channelid": *channelID}).One(&result)
	return &result, error
}

// GetAutoMessages returns all automessages for specific channel
func GetAutoMessages(channelID *string) ([]models.AutoMessageWithHistory, error) {
	var result []models.AutoMessageWithHistory
	error := db.C(autoMessageCollectionName).Find(bson.M{"channelid": *channelID}).All(&result)
	return result, error
}

// CreateAutoMessage inserts new automessage to channel, specified in AutoMessageUpdate object
func CreateAutoMessage(autoMessageUpdate *models.AutoMessageUpdate) (*bson.ObjectId, error) {
	id := bson.NewObjectId()
	now := time.Now()
	if strings.TrimSpace(autoMessageUpdate.Message) == "" || autoMessageUpdate.DurationLimit < 60 || autoMessageUpdate.MessageLimit < 20 {
		return nil, errors.New("Validation Failed")
	}
	var durationLimit = time.Second * time.Duration(autoMessageUpdate.DurationLimit)
	db.C(autoMessageCollectionName).Insert(
		models.AutoMessageWithHistory{
			AutoMessage: models.AutoMessage{
				ID:                id,
				ChannelID:         autoMessageUpdate.ChannelID,
				Message:           autoMessageUpdate.Message,
				MessageThreshold:  autoMessageUpdate.MessageLimit,
				MessageLimit:      autoMessageUpdate.MessageLimit,
				Game:              autoMessageUpdate.Game,
				DurationLimit:     durationLimit,
				DurationThreshold: now.Add(durationLimit)},
			History: []models.AutoMessageHistory{
				models.AutoMessageHistory{
					User:          autoMessageUpdate.User,
					UserID:        autoMessageUpdate.UserID,
					Game:          autoMessageUpdate.Game,
					Date:          now,
					Message:       autoMessageUpdate.Message,
					MessageLimit:  autoMessageUpdate.MessageLimit,
					DurationLimit: durationLimit}}})
	return &id, nil
}

// UpdateAutoMessage updates specific automessage, which details are specified in AutoMessageUpdate object
func UpdateAutoMessage(autoMessageUpdate *models.AutoMessageUpdate) error {
	if autoMessageUpdate.DurationLimit < 60 || autoMessageUpdate.MessageLimit < 20 {
		return errors.New("Validation Failed")
	}
	now := time.Now()
	var durationLimit = time.Second * time.Duration(autoMessageUpdate.DurationLimit)
	db.C(autoMessageCollectionName).Update(
		bson.M{"_id": bson.ObjectIdHex(autoMessageUpdate.ID), "channelid": autoMessageUpdate.ChannelID},
		bson.M{
			"$push": bson.M{
				"history": bson.M{
					"$each": []models.AutoMessageHistory{models.AutoMessageHistory{
						User:          autoMessageUpdate.User,
						UserID:        autoMessageUpdate.UserID,
						Game:          autoMessageUpdate.Game,
						Date:          now,
						Message:       autoMessageUpdate.Message,
						MessageLimit:  autoMessageUpdate.MessageLimit,
						DurationLimit: durationLimit}},
					"$sort":  bson.M{"date": -1},
					"$slice": 5}},

			"$set": models.AutoMessage{
				ChannelID:         autoMessageUpdate.ChannelID,
				Message:           autoMessageUpdate.Message,
				MessageThreshold:  autoMessageUpdate.MessageLimit,
				MessageLimit:      autoMessageUpdate.MessageLimit,
				Game:              autoMessageUpdate.Game,
				DurationLimit:     durationLimit,
				DurationThreshold: now.Add(durationLimit)}})
	return nil
}
