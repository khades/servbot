package models

import (
	"time"
)

type SubTrain struct {
	Enabled              bool      `json:"enabled"`
	ExpirationLimit      int       `json:"expirationLimit"`
	NotificationLimit    int       `json:"notificationLimit"`
	NotificationShown    bool      `json:"notificationShown"`
	ExpirationTime       time.Time `json:"expirationTime"`
	NotificationTime     time.Time `json:"notificationTime"`
	AppendTemplate       string    `json:"appendTemplate"`
	TimeoutTemplate      string    `json:"timeoutTemplate"`
	NotificationTemplate string    `json:"notificationTemplate"`
	CurrentStreak        int       `json:"сurrentStreak"`
}
