package repos

import (
	"math/rand"
	"time"

	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
	"github.com/globalsign/mgo/bson"
)

var subdayCollection = "subdays"

// GetLastActiveSubday returns last active subday for specified channel
func GetLastActiveSubday(channelID *string) (*models.Subday, error) {
	var result models.Subday
	error := db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
		"isactive":  true}).One(&result)
	return &result, error
}

// GetLastSubday returns any last subday for specified channel
func GetLastSubday(channelID *string) (*models.SubdayNoWinners, error) {
	var result models.SubdayNoWinners
	error := db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
	}).Sort("-date").One(&result)
	return &result, error
}

// GetLastSubdayMod returns any last subday for specified channel with extended information about it
func GetLastSubdayMod(channelID *string) (*models.Subday, error) {
	var result models.Subday
	error := db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
	}).Sort("-date").One(&result)
	return &result, error
}

// GetSubdayByID returns specified subday
func GetSubdayByID(id *string) (*models.SubdayNoWinners, error) {
	var result models.SubdayNoWinners
	error := db.C(subdayCollection).Find(bson.M{
		"_id": bson.ObjectIdHex(*id)}).One(&result)
	return &result, error
}

// GetSubdayByIDMod returns specified subday with extended information about it
func GetSubdayByIDMod(id *string) (*models.Subday, error) {
	var result models.Subday
	error := db.C(subdayCollection).Find(bson.M{
		"_id": bson.ObjectIdHex(*id)}).One(&result)
	return &result, error
}

// GetSubdays return all subdays for specified channel
func GetSubdays(channelID *string) ([]models.SubdayList, error) {
	var result []models.SubdayList
	error := db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
	}).Sort("-date").All(&result)
	return result, error
}

// CloseActiveSubday closes any active subday on specified channel
func CloseActiveSubday(channelID *string) {
	db.C(subdayCollection).UpdateAll(bson.M{
		"channelid": *channelID,
		"isactive":  true},
		bson.M{"$set": bson.M{"isactive": false}})
	SetSubdayIsActive(channelID, false)

}

// CloseSubday closes specified subday on specified channel
func CloseSubday(channelID *string, id *string) {
	db.C(subdayCollection).UpdateAll(bson.M{
		"channelid": *channelID,
		"_id":       bson.ObjectIdHex(*id)},
		bson.M{"$set": bson.M{"isactive": false}})
	SetSubdayIsActive(channelID, false)

}

// SubdayPullWinner pulls specified user from winners in specified subday
func SubdayPullWinner(id *string, user *string) {
	db.C(subdayCollection).UpdateAll(bson.M{
		"_id":      bson.ObjectIdHex(*id),
		"isactive": true},
		bson.M{
			"$pull": bson.M{
				"winners": bson.M{"user": *user}}})
}

// PushWinners pushes (replaces) winners list for specified subday
func PushWinners(id bson.ObjectId, winners []models.SubdayRecord) {
	db.C(subdayCollection).UpdateAll(bson.M{
		"_id":      id,
		"isactive": true}, bson.M{
		"$set": bson.M{"winners": winners},
		"$push": bson.M{"winnershistory": bson.M{
			"$each": []models.SubdayWinnersHistory{models.SubdayWinnersHistory{
				Winners: winners,
				Date:    time.Now()}},
			"$sort":  bson.M{"date": -1},
			"$slice": 5}}})
}

// PickRandomWinnerForSubday rolls one winner and automatically upserts him to winners in specified subday for specified channel (this was made to enforce that only streamer could roll)
func PickRandomWinnerForSubday(channelID *string, id *string) *models.SubdayRecord {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "subday",
		"action":  "PickRandomWinnerForSubday"})
	var result models.Subday
	error := db.C(subdayCollection).Find(bson.M{"_id": bson.ObjectIdHex(*id), "channelid": *channelID, "isactive": true}).One(&result)
	if error != nil {
		logger.Debugf("some error with collection")
		return nil
	}
	votes := []models.SubdayRecord{}
	if len(result.Votes) == 0 {
		logger.Debugf("No Votes")

		return nil
	}
	if len(result.Votes) == len(result.Winners) {
		logger.Debugf("Nobody To Pick")

		return nil
	}
	winners := result.Winners
	for _, vote := range result.Votes {
		found := false
		for _, winner := range result.Winners {
			if vote.UserID == winner.UserID {
				found = true
				break
			}
		}
		if found == false {
			votes = append(votes, vote)
		}
	}
	var winner models.SubdayRecord
	random := rand.Intn(len(votes))
	winner = votes[random]
	winners = append(winners, winner)
	PushWinners(bson.ObjectIdHex(*id), winners)
	return &winner
}

// VoteForSubday upserts specified vote variant of specified user for specified subday
func VoteForSubday(user *string, userID *string, id *bson.ObjectId, game *string) {

	db.C(subdayCollection).UpdateAll(bson.M{
		"_id":          *id,
		"votes.userid": userID},
		bson.M{"$set": bson.M{
			"votes.$.game": *game,
			"votes.$.user": *user}})

	db.C(subdayCollection).UpdateAll(bson.M{
		"_id":          *id,
		"votes.userid": bson.M{"$ne": userID}},
		bson.M{"$push": bson.M{"votes": models.SubdayRecord{
			User:   *user,
			UserID: *userID,
			Game:   *game}}})

}

// CreateNewSubday creates new subday IF there's no active subdays
func CreateNewSubday(channelID *string, subsOnly bool, name *string) bool {
	_, subdayError := GetLastActiveSubday(channelID)
	if subdayError == nil {
		return false
	}
	object := models.Subday{Name: *name, SubsOnly: subsOnly, IsActive: true, Date: time.Now(), ChannelID: *channelID}
	db.C(subdayCollection).Insert(object)
	SetSubdayIsActive(channelID, true)
	return true
}
