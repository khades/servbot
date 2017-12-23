package repos

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"github.com/khades/servbot/models"
)

var subdayCollection string = "subdays"

func GetLastActiveSubday(channelID *string) (*models.Subday, error) {
	var result models.Subday
	error := Db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
		"isactive": true}).One(&result)
	return &result, error
}

func GetLastSubday(channelID *string) (*models.Subday, error) {
	var result models.Subday
	error := Db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
		}).Sort("-date").One(&result)
	return &result, error
}

func CloseActiveSubday(channelID *string) {
	Db.C(subdayCollection).UpdateAll(bson.M{
		"channelid": *channelID,
		"isactive": true},
		 bson.M{"isactive":false})
}

func PullWinner(id *bson.ObjectId, user *string) {
	Db.C(subdayCollection).UpdateAll(bson.M{
		"_id": *id},
		 bson.M{
			 "$pull": bson.M{
				 "winners": bson.M{"user":*user}}})
}
func PushWinners(id *bson.ObjectId, winners *[]models.SubdayRecord) {
	Db.C(subdayCollection).UpdateAll(bson.M{"_id": *id}, bson.M{
		"$set": bson.M{"winners": *winners},
		"$push": bson.M{"winnershistory":
			bson.M{
				"$each":  []models.SubdayWinnersHistory{models.SubdayWinnersHistory{
					Winners: *winners,
					Date: time.Now()}},
				"$sort":  bson.M{"date": -1},
				"$slice": 5}}})
}
func PickRandomWinnerForSubday(id *bson.ObjectId, channelID *string) bool {
	var result models.Subday
	error := Db.C(subdayCollection).Find(bson.M{"_id": *id, "channelid": *channelID, "isactive":true}).One(&result)
	if error !=nil {
		return false
	}
	votes := []models.SubdayRecord{}
	if (len(votes) == 0) {
		return false
	}
	winners := result.Winners
	for _, vote := range result.Votes {
		found := false
		for _, winner := range result.Winners {
			if (vote.UserID == winner.UserID) {
				found = true
				break
			}
		}
		if (found == false){
			votes = append(votes, vote)
		}
	}
	winners = append(winners, votes[rand.Intn(len(votes) -1)])
	return true
}
func VoteForSubday(user *string, userID *string, id *bson.ObjectId, game *string) {
	
	Db.C(subdayCollection).UpdateAll(bson.M{
		"_id" : *id,
		"votes.userid": userID},
		bson.M{"$set":bson.M{"votes": bson.M{
			"votes.*.game": *game,
			"votes.*.user": *user}}})
			
	Db.C(subdayCollection).UpdateAll(bson.M{
		"_id": *id,
		"votes.userid": bson.M{"$ne":userID}},
		 bson.M{"$push":bson.M{"votes": models.SubdayRecord{
			 User: *user,
			 UserID: *userID,
			 Game: *game}}})

}

func CreateNewSubday(channelID *string, subsOnly bool, name *string) bool {
	_, subdayError := GetLastActiveSubday(channelID)
	if subdayError != nil {
		return false
	}
	Db.C(subdayCollection).Insert(models.Subday{Name: *name, SubsOnly: subsOnly, Date: time.Now(), ChannelID:*channelID})
	return true
}