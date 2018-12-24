package subscriptionInfo

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// SubscriptionInfo shows how many t imes user is subscribed to a channel
type SubscriptionInfo struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Count     int           `json:"count"`
	IsPrime   bool          `json:"isPrime"`
	User      string        `json:"user"`
	UserID    string        `json:"userID"`
	ChannelID string        `json:"channelID"`
	Date      time.Time     `json:"date"`
	SubPlan   string        `json:"subPlan"`
}
