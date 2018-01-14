package repos

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var autoMessageCollectionName = "autoMessages"

func DecrementAutoMessages(channelID *string) {
	channelInfo, error := GetChannelInfo(channelID)
	games := []string{""}
	if error == nil && channelInfo.StreamStatus.Online == true {
		games = append(games, channelInfo.StreamStatus.Game)
	}
	Db.C(autoMessageCollectionName).UpdateAll(bson.M{
		"channelid": *channelID,
		"message":   bson.M{"$ne": ""},
		"$or": []bson.M{
			bson.M{"game": bson.M{"$in": games}},
			bson.M{"game": bson.M{"$exists": false}}}},
		bson.M{"$inc": bson.M{"messagethreshold": -1}})
}
func RemoveInactiveAutoMessages(channelID *string) (*[]models.AutoMessage, error) {
	var result []models.AutoMessage
	error := Db.C(autoMessageCollectionName).Find(bson.M{
		"channelid": *channelID,
		"message":   "",
		"history.date": 
			bson.M{"$not": bson.M{"$gte": time.Now().Add(24 * -7 * time.Hour)}}}).All(&result)
	return &result, error
}
func GetCurrentAutoMessages() (*[]models.AutoMessage, error) {
	//log.Println("AutoMessage: Getting Current AutoMessages")
	var result []models.AutoMessage
	error := Db.C(autoMessageCollectionName).Find(bson.M{
		"message":           bson.M{"$ne": ""},
		"messagethreshold":  bson.M{"$lte": 0},
		"durationthreshold": bson.M{"$lte": time.Now()}}).All(&result)
	log.Printf("AutoMessage: Got %d AutoMessages", len(result))
	//log.Println(error)
	return &result, error
}

func ResetAutoMessageThreshold(autoMessage *models.AutoMessage) {
	log.Printf("AutoMessage: Resetting AutoMessage %s", autoMessage.ID)
	now := time.Now()
	Db.C(autoMessageCollectionName).Update(bson.M{"_id": autoMessage.ID}, bson.M{"$set": bson.M{
		"messagethreshold":  autoMessage.MessageLimit,
		"durationthreshold": now.Add(autoMessage.DurationLimit)}})
}

func GetAutoMessage(id *string, channelID *string) (*models.AutoMessageWithHistory, error) {
	var result models.AutoMessageWithHistory
	objectID := bson.ObjectIdHex(*id)
	error := Db.C(autoMessageCollectionName).Find(bson.M{"_id": objectID, "channelid": *channelID}).One(&result)
	return &result, error
}

func GetAutoMessages(channelID *string) (*[]models.AutoMessageWithHistory, error) {
	var result []models.AutoMessageWithHistory
	error := Db.C(autoMessageCollectionName).Find(bson.M{"channelid": *channelID}).All(&result)
	return &result, error
}
func CreateAutoMessage(autoMessageUpdate *models.AutoMessageUpdate) (*bson.ObjectId, error) {
	id := bson.NewObjectId()
	now := time.Now()
	if strings.TrimSpace(autoMessageUpdate.Message) == "" || autoMessageUpdate.DurationLimit < 60 || autoMessageUpdate.MessageLimit < 20 {
		return nil, errors.New("Validation Failed")
	}
	var durationLimit = time.Second * time.Duration(autoMessageUpdate.DurationLimit)
	Db.C(autoMessageCollectionName).Insert(
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

func UpdateAutoMessage(autoMessageUpdate *models.AutoMessageUpdate) error {
	if autoMessageUpdate.DurationLimit < 60 || autoMessageUpdate.MessageLimit < 20 {
		return errors.New("Validation Failed")
	}
	now := time.Now()
	var durationLimit = time.Second * time.Duration(autoMessageUpdate.DurationLimit)
	Db.C(autoMessageCollectionName).Update(
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
