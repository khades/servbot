package repos

import (
	"log"

	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var bitsCollection = "bits"

func AddBitsToUser(channelID *string, userID *string, user *string, amount int, reason string) {
	log.Println(*channelID)
	log.Println(*userID)
	log.Println(amount)

	Db.C(bitsCollection).Upsert(bson.M{
		"channelid": *channelID,
		"userid":    *userID},
		bson.M{
			"$inc": bson.M{"amount": amount},
			"$set": bson.M{"user": *user},
			"$push": bson.M{
				"history": bson.M{
					"$each":  []models.UserBitsHistory{models.UserBitsHistory{Reason: reason, Change: amount}},
					"$sort":  bson.M{"date": -1},
					"$slice": 100}}})
}

func GetBitsForChannel(channelID *string) (*[]models.UserBits, error) {
	var result []models.UserBits
	error := Db.C(bitsCollection).Find(models.ChannelSelector{ChannelID: *channelID}).All(&result)
	return &result, error
}

func GetBitsForChannelUser(channelID *string, userID *string) (*models.UserBitsWithHistory, error) {
	var result models.UserBitsWithHistory
	error := Db.C(bitsCollection).Find(bson.M{
		"channelid": *channelID,
		"userid":    *userID}).One(&result)
	return &result, error
}

func PutSubscriptionBits(channelID *string, userID *string, user *string, subPlan *string) {
	switch *subPlan {
	case "Prime":
		{
			AddBitsToUser(channelID, userID, user, 499, "subprime")
		}
	case "1000":
		{
			AddBitsToUser(channelID, userID, user, 499, "sub")
		}
	case "2000":
		{
			AddBitsToUser(channelID, userID, user, 999, "sub")
		}
	case "3000":
		{
			AddBitsToUser(channelID, userID, user, 2499, "sub")
		}
	}
}
