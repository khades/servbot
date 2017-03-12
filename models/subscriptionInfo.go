package models

import (
	"time"
)

// SubscriptionInfo shows how many times user is subscribed to a channel
type SubscriptionInfo struct {
	Count     int       `json:"count"`
	IsPrime   bool      `json:"isPrime"`
	User      string    `json:"user"`
	UserID    string    `json:"userID"`
	ChannelID string    `json:"channelID"`
	Date      time.Time `json:"date"`
}
