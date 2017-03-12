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

func GetSubsForChannel(channelID *string, limit time.Time) (*[]models.SubscriptionInfo, error) {
	var result []models.SubscriptionInfo
	// var localLimit time.Time
	// if limit.IsZero() {
	// 	day := -24 * time.Hour
	// 	localLimit = time.Now().Add(day)
	// } else {
	// 	localLimit = limit
	// }
	error := Db.C(subscriptionInfoCollection).Find(bson.M{"channelid": *channelID}).Sort("-date").All(&result)
	return &result, error
}
