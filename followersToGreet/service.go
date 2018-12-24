package followersToGreet

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	collection *mgo.Collection
}

func (service *Service) Add(channelID *string, follower string) {
	service.collection.Upsert(bson.M{"channelid": *channelID}, bson.M{"$addToSet":bson.M{"followers":follower}})
}

func (service *Service) Reset(channelID *string) {
	service.collection.Remove(bson.M{"channelid": *channelID})
}

func (service *Service) List() ([]FollowersToGreet,error) {
	result := []FollowersToGreet{}
	err := service.collection.Find(bson.M{}).All(&result)
	return result, err
}