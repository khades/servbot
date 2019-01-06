package balance

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	collection *mgo.Collection
}

func (service *Service) Inc(channelID string, userID string, user string, amount float64) {
	if user != "" {
		service.collection.Upsert(
			bson.M{"channelid": channelID, "userid": userID},
			bson.M{
				"$inc": bson.M{"balance": amount},
				"$set": bson.M{"user": user}})
	} else {
		service.collection.Upsert(
			bson.M{"channelid": channelID, "userid": userID},
			bson.M{
				"$inc": bson.M{"balance": amount}})
	}
}

func (service *Service) Dec(channelID string, userID string, user string, amount float64) bool {
	if user != "" {
		err := service.collection.Update(
			bson.M{"channelid": channelID, "userid": userID, "balance": bson.M{"$gte": amount}},
			bson.M{
				"$inc": bson.M{"balance": -amount},
				"$set": bson.M{"user": user}})
		if err == nil {
			return true
		} else {
			return false
		}
	} else {
		err := service.collection.Update(
			bson.M{"channelid": channelID, "userid": userID, "balance": bson.M{"$gte": amount}},
			bson.M{
				"$inc": bson.M{"balance": -amount}})
		if err == nil {
			return true
		} else {
			return false
		}
	}
}

func (service *Service) Get(channelID string, userID string) (Balance, error) {
	result := Balance{}
	error := service.collection.Find(bson.M{"channelid": channelID, "userid": userID}).One(&result)
	return result, error
}

func (service *Service) List(channelID string) ([]Balance, error) {
	result := []Balance{}
	error := service.collection.Find(bson.M{"channelid": channelID}).All(&result)
	return result, error
}
