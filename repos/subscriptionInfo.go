package repos

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

var subscriptionInfoCollection string = "subsciptionInfo"

func LogSubscription(info *models.SubscriptionInfo) {
	Db.C(subscriptionInfoCollection).Insert(*info)
}
func GetSubsForChannelWithLimit(channelID *string, limit time.Time) (*[]models.SubscriptionInfo, error) {
	var result []models.SubscriptionInfo
	error := Db.C(subscriptionInfoCollection).Find(bson.M{
		"channelid": *channelID,
		"date":      bson.M{"$gte": limit}}).Sort("-date").All(&result)
	return &result, error
}

func GetSubsForChannel(channelID *string) (*[]models.SubscriptionInfo, error) {
	var result []models.SubscriptionInfo
	day := -24 * time.Hour
	localLimit := time.Now().Add(day)
	error := Db.C(subscriptionInfoCollection).Find(bson.M{
		"channelid": *channelID,
		"date":      bson.M{"$gte": localLimit}}).Sort("-date").All(&result)
	return &result, error
}
