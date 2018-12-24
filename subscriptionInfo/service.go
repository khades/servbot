package subscriptionInfo

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)


type Service struct {
	collection *mgo.Collection
}

// Log writes user subscription
func (service *Service) Log(info *SubscriptionInfo) {
	//db.C(subscriptionInfoCollection).Insert(*info)
	service.collection.Upsert(bson.M{"userid": info.UserID, "channelid": info.ChannelID}, info)
	//PutSubscriptionBits(&info.ChannelID, &info.UserID, &info.User, &info.SubPlan)
}

// Get returns list of subscription for specified channel after specified time
func (service *Service) Get(channelID *string, limit time.Time) ([]SubscriptionInfo, error) {
	var result []SubscriptionInfo
	error := service.collection.Find(bson.M{
		"channelid": *channelID,
		"date":      bson.M{"$gte": limit}}).Sort("-date").All(&result)
	return result, error
}

// GetDefault is version of Get with pre-built time of three days
func (service *Service) GetDefault(channelID *string) ([]SubscriptionInfo, error) {
	day := -24 * 3 * time.Hour
	localLimit := time.Now().Add(day)
	return service.Get(channelID, localLimit)
}
