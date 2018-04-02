package repos

import (
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/khades/servbot/models"
)

var subscriptionInfoCollection  = "subsciptionInfo"

// LogSubscription writes user subscription 
func LogSubscription(info *models.SubscriptionInfo) {
	db.C(subscriptionInfoCollection).Insert(*info)
	//PutSubscriptionBits(&info.ChannelID, &info.UserID, &info.User, &info.SubPlan)
}

// GetSubsForChannelWithLimit returns list of subscription for specified channel after specified time 
func GetSubsForChannelWithLimit(channelID *string, limit time.Time) ([]models.SubscriptionInfo, error) {
	var result []models.SubscriptionInfo
	error := db.C(subscriptionInfoCollection).Find(bson.M{
		"channelid": *channelID,
		"date":      bson.M{"$gte": limit}}).Sort("-date").All(&result)
	return result, error
}

// GetSubsForChannel is version of GetSubsForChannelWithLimit with pre-built time of three days
func GetSubsForChannel(channelID *string) ([]models.SubscriptionInfo, error) {
	day := -24 * 3 * time.Hour
	localLimit := time.Now().Add(day)
	return GetSubsForChannelWithLimit(channelID, localLimit)
}
