package subscriptionInfoAPI

import (
	"github.com/khades/servbot/subscriptionInfo"
	"time"
)

type subscriptionEvent struct {
	Subscription     subscriptionInfo.SubscriptionInfo `json:"subscription"`
	CurrentCallTime  time.Time               `json:"currentCallTimetime"`
	PreviousCallTime time.Time               `json:"previousCallTimetime"`
}