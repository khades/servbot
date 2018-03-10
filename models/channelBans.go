package models

import (
	"time"
)

// ChannelBanRecord struct describes one specific ban 
type ChannelBanRecord struct {
	User string `json:"user"`
	UserID string `json:"userID"`
	BanLength int `json:"banLength"`
	Date time.Time `json:"date"`
}

// ChannelBans struct descibes last bans on specified channel
type ChannelBans struct {
	ChannelID string `json:"channelID"`
	Bans []ChannelBanRecord `json:"bans"`
}