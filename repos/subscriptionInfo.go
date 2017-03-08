package repos

import "github.com/khades/servbot/models"

var subscriptionInfoCollection string = "subsciptionInfo"

func LogSubscription(info *models.SubscriptionInfo) {
	Db.C(subscriptionInfoCollection).Insert(*info)
}
