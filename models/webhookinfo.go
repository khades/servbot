package models

import "time"

type WebHookInfo struct {
	ChannelID string
	Topic     string
	Secret    string
	ExpiresAt time.Time
}
