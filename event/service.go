package event

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	collection *mgo.Collection
}

func (service *Service) Put(channelID string, event Event) {
	event.Date = time.Now()
	service.collection.Update(bson.M{"channelid": channelID}, bson.M{
		"$push": bson.M{
			"events": bson.M{"$each": []Event{event},
				"$sort":  bson.M{"date": -1},
				"$slice": 100},
		},
		"$pull": bson.M{"events": bson.M{
			"date": bson.M{"$lte": time.Now().Add(-3 * 24 * time.Hour)}}},
	})
}

func (service *Service) SetViewedTill(channelID string, viewedTill time.Time) {
	service.collection.Update(
		bson.M{"channelid": channelID},
		bson.M{"$set": bson.M{
			"viewedtill": viewedTill,
		}})
}

func (service *Service) Get(channelID string) (*Events, error) {
	result := &Events{}
	err := service.collection.Find(bson.M{"channelid": channelID}).One(result)
	return result, err
}
