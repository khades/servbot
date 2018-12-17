package models

import (
	"time"
)

// SubTrain descibes subtrain settings and state on specified channel
type SubTrain struct {
	Enabled              bool      `json:"enabled"`
	OnlyNewSubs          bool      `json:"onlyNewSubs"`
	ExpirationLimit      int       `json:"expirationLimit"`
	NotificationLimit    int       `json:"notificationLimit"`
	NotificationShown    bool      `json:"notificationShown"`
	ExpirationTime       time.Time `json:"expirationTime"`
	NotificationTime     time.Time `json:"notificationTime"`
	AppendTemplate       string    `json:"appendTemplate"`
	TimeoutTemplate      string    `json:"timeoutTemplate"`
	NotificationTemplate string    `json:"notificationTemplate"`
	CurrentStreak        int       `json:"currentStreak"`
	Users                []string  `json:"users"`
}
