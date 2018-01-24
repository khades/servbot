package repos

import (
	"time"
	"log"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"github.com/khades/servbot/models"
)

var subdayCollection  = "subdays"

func GetLastActiveSubday(channelID *string) (*models.Subday, error) {
	var result models.Subday
	error := Db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
		"isactive": true}).One(&result)
	return &result, error
}

func GetLastSubday(channelID *string) (*models.SubdayNoWinners, error) {
	var result models.SubdayNoWinners
	error := Db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
		}).Sort("-date").One(&result)
	return &result, error
}

func GetLastSubdayMod(channelID *string) (*models.Subday, error) {
	var result models.Subday
	error := Db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
		}).Sort("-date").One(&result)
	return &result, error
}

func GetSubdayById(id *string) (*models.SubdayNoWinners, error) {
	var result models.SubdayNoWinners
	error := Db.C(subdayCollection).Find(bson.M{
		"_id": bson.ObjectIdHex(*id)}).One(&result)
	return &result, error
}

func GetSubdayByIdMod(id *string) (*models.Subday, error) {
	var result models.Subday
	error := Db.C(subdayCollection).Find(bson.M{
		"_id": bson.ObjectIdHex(*id)}).One(&result)
	return &result, error
}

func GetSubdays(channelID *string)(*[]models.SubdayList, error) {
	var result []models.SubdayList
	error := Db.C(subdayCollection).Find(bson.M{
		"channelid": *channelID,
		}).Sort("-date").All(&result)
	return &result, error
}

func CloseActiveSubday(channelID *string) {
	Db.C(subdayCollection).UpdateAll(bson.M{
		"channelid": *channelID,
		"isactive": true},
		 bson.M{"$set":bson.M{"isactive":false}})
	SetSubdayIsActive(channelID, false)

}

func CloseSubday(channelID *string, id *string) {
	Db.C(subdayCollection).UpdateAll(bson.M{
		"channelid": *channelID,
		"_id": bson.ObjectIdHex(*id)},
		 bson.M{"$set":bson.M{"isactive":false}})
	SetSubdayIsActive(channelID, false)

}
func SubdayPullWinner(id *string, user *string) {
	Db.C(subdayCollection).UpdateAll(bson.M{
		"_id": bson.ObjectIdHex(*id),
		"isactive": true},
		 bson.M{
			 "$pull": bson.M{
				 "winners": bson.M{"user":*user}}})
}

func PushWinners(id bson.ObjectId, winners *[]models.SubdayRecord) {
	Db.C(subdayCollection).UpdateAll(bson.M{
		"_id":  id,
		"isactive": true}, bson.M{
		"$set": bson.M{"winners": *winners},
		"$push": bson.M{"winnershistory":
			bson.M{
				"$each":  []models.SubdayWinnersHistory{models.SubdayWinnersHistory{
					Winners: *winners,
					Date: time.Now()}},
				"$sort":  bson.M{"date": -1},
				"$slice": 5}}})
}
func PickRandomWinnerForSubday( channelID *string,id *string) *models.SubdayRecord {
	var result models.Subday
	error := Db.C(subdayCollection).Find(bson.M{"_id": bson.ObjectIdHex(*id), "channelid": *channelID, "isactive":true}).One(&result)
	if error !=nil {
		log.Println("some error with collection")
		return nil
	}
	votes := []models.SubdayRecord{}
	if (len(result.Votes) == 0) {
		log.Println("No Votes")

		return nil
	}
	if (len(result.Votes) == len(result.Winners)) {
		log.Println("Nobody To Pick")

		return nil
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
	var winner models.SubdayRecord
	random := rand.Intn(len(votes))
	winner = votes[random]
	winners = append(winners,winner)
	PushWinners(bson.ObjectIdHex(*id), &winners)
	return &winner
}

func VoteForSubday(user *string, userID *string, id *bson.ObjectId, game *string) {

	Db.C(subdayCollection).UpdateAll(bson.M{
		"_id" : *id,
		"votes.userid": userID},
		bson.M{"$set":bson.M{
			"votes.$.game": *game,
			"votes.$.user": *user}})
			
	Db.C(subdayCollection).UpdateAll(bson.M{
		"_id": *id,
		"votes.userid": bson.M{"$ne":userID}},
		 bson.M{"$push":bson.M{"votes": models.SubdayRecord{
			 User: *user,
			 UserID: *userID,
			 Game: *game}}})

}


func SetSubdayIsActive(channelID *string, isActive bool) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.SubdayIsActive= isActive
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, SubdayIsActive: isActive})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subdayisactive": isActive}})
}

func CreateNewSubday(channelID *string, subsOnly bool, name *string) bool {
	_, subdayError := GetLastActiveSubday(channelID)
	if subdayError == nil {
		return false
	}
	object :=models.Subday{Name: *name, SubsOnly: subsOnly, IsActive: true, Date: time.Now(), ChannelID:*channelID}
	Db.C(subdayCollection).Insert(object)
	SetSubdayIsActive(channelID, true)
	return true
}