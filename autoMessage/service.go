package autoMessage

import (
	"errors"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/channelInfo"
	"github.com/sirupsen/logrus"
)

type Service struct {
	// Own Fields
	collection *mgo.Collection
}

// Decrement decrement message threshold every time someone writes in chat, with check if autoMessage is bound to specific game on stream
func (service *Service) Decrement(channelInfo *channelInfo.ChannelInfo) {
	games := []string{""}
	if channelInfo.StreamStatus.Online == true {
		games = append(games, channelInfo.StreamStatus.Game)
	}
	service.collection.UpdateAll(bson.M{
		"channelid": channelInfo.ChannelID,
		"message":   bson.M{"$ne": ""},
		"$or": []bson.M{
			bson.M{"game": bson.M{"$in": games}},
			bson.M{"game": bson.M{"$exists": false}}}},
		bson.M{"$inc": bson.M{"messagethreshold": -1}})
}

// RemoveInactive removes all automessages on channel, which had no update for a week
func (service *Service) RemoveInactive(channelID *string) {
	service.collection.RemoveAll(bson.M{
		"channelid":    *channelID,
		"message":      "",
		"history.date": bson.M{"$not": bson.M{"$gte": time.Now().Add(24 * -7 * time.Hour)}}})
}

// ListCurrent returns list of ALL autoMessage, which are served on that chatbot instance
func (service *Service) ListCurrent() ([]AutoMessage, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "autoMessage",
		"action":  "ListCurrent"})
	logger.Debug("AutoMessage: Getting Current AutoMessages")
	var result []AutoMessage
	error := service.collection.Find(bson.M{
		"message":           bson.M{"$ne": ""},
		"messagethreshold":  bson.M{"$lte": 0},
		"durationthreshold": bson.M{"$lte": time.Now()}}).All(&result)
	logger.Infof("AutoMessage: Got %d AutoMessages", len(result))
	return result, error
}

// ResetThreshold resets autoMessage thresholds after successfull autoMessage execution
func (service *Service) ResetThreshold(autoMessage *AutoMessage) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "autoMessage",
		"action":  "ResetThreshold"})
	logger.Infof("AutoMessage: Resetting AutoMessage %s", autoMessage.ID)
	now := time.Now()
	service.collection.Update(bson.M{"_id": autoMessage.ID}, bson.M{"$set": bson.M{
		"messagethreshold":  autoMessage.MessageLimit,
		"durationthreshold": now.Add(autoMessage.DurationLimit)}})
}

// Search returns specific autoMessage for specific channel
func (service *Service) Get(id *string, channelID *string) (*AutoMessageWithHistory, error) {
	var result AutoMessageWithHistory
	objectID := bson.ObjectIdHex(*id)
	error := service.collection.Find(bson.M{"_id": objectID, "channelid": *channelID}).One(&result)
	return &result, error
}

// List returns all automessages for specific channel
func (service *Service) List(channelID *string) ([]AutoMessageWithHistory, error) {
	var result []AutoMessageWithHistory
	error := service.collection.Find(bson.M{"channelid": *channelID}).All(&result)
	return result, error
}

// Create inserts new autoMessage to channel, specified in AutoMessageUpdate object
func (service *Service) Create(autoMessageUpdate *AutoMessageUpdate) (*bson.ObjectId, error) {
	id := bson.NewObjectId()
	now := time.Now()
	if strings.TrimSpace(autoMessageUpdate.Message) == "" || autoMessageUpdate.DurationLimit < 60 || autoMessageUpdate.MessageLimit < 20 {
		return nil, errors.New("Validation Failed")
	}
	var durationLimit = time.Second * time.Duration(autoMessageUpdate.DurationLimit)
	service.collection.Insert(
		AutoMessageWithHistory{
			AutoMessage: AutoMessage{
				ID:                id,
				ChannelID:         autoMessageUpdate.ChannelID,
				Message:           autoMessageUpdate.Message,
				MessageThreshold:  autoMessageUpdate.MessageLimit,
				MessageLimit:      autoMessageUpdate.MessageLimit,
				Game:              autoMessageUpdate.Game,
				DurationLimit:     durationLimit,
				DurationThreshold: now.Add(durationLimit)},
			History: []AutoMessageHistory{
				AutoMessageHistory{
					User:          autoMessageUpdate.User,
					UserID:        autoMessageUpdate.UserID,
					Game:          autoMessageUpdate.Game,
					Date:          now,
					Message:       autoMessageUpdate.Message,
					MessageLimit:  autoMessageUpdate.MessageLimit,
					DurationLimit: durationLimit}}})
	return &id, nil
}

// Update updates specific autoMessage, which details are specified in AutoMessageUpdate object
func (service *Service) Update(autoMessageUpdate *AutoMessageUpdate) error {
	if autoMessageUpdate.DurationLimit < 60 || autoMessageUpdate.MessageLimit < 20 {
		return errors.New("Validation Failed")
	}
	now := time.Now()
	var durationLimit = time.Second * time.Duration(autoMessageUpdate.DurationLimit)
	service.collection.Update(
		bson.M{"_id": bson.ObjectIdHex(autoMessageUpdate.ID), "channelid": autoMessageUpdate.ChannelID},
		bson.M{
			"$push": bson.M{
				"history": bson.M{
					"$each": []AutoMessageHistory{AutoMessageHistory{
						User:          autoMessageUpdate.User,
						UserID:        autoMessageUpdate.UserID,
						Game:          autoMessageUpdate.Game,
						Date:          now,
						Message:       autoMessageUpdate.Message,
						MessageLimit:  autoMessageUpdate.MessageLimit,
						DurationLimit: durationLimit}},
					"$sort":  bson.M{"date": -1},
					"$slice": 5}},

			"$set": AutoMessage{
				ChannelID:         autoMessageUpdate.ChannelID,
				Message:           autoMessageUpdate.Message,
				MessageThreshold:  autoMessageUpdate.MessageLimit,
				MessageLimit:      autoMessageUpdate.MessageLimit,
				Game:              autoMessageUpdate.Game,
				DurationLimit:     durationLimit,
				DurationThreshold: now.Add(durationLimit)}})
	return nil
}
