package models

import (
	"time"
)

// SubscriptionInfo shows how many times user is subscribed to a channel
type SubscriptionInfo struct {
	Count     int
	IsPrime   bool
	User      string
	UserID    string
	ChannelID string
	Date      time.Time
}
