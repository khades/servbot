package models

import "time"

type WebHookInfo struct {
	ChannelID string
	Topic     string
	Challenge string
	ExpiresAt time.Time
}
