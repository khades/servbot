package models

import (
	"time"
)

type ChannelBanRecord struct {
	User string `json:"user"`
	UserID string `json:"userID"`
	BanLength int `json:"banLength"`
	Date time.Time `json:"date"`
}
type ChannelBans struct {
	ChannelID string `json:"channelID"`
	Bans []ChannelBanRecord `json:"bans"`
}