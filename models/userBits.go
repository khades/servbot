package models

import "time"

// UserBits is obsolete
type UserBits struct {
	User      string `json:"user"`
	UserID    string `json:"userID"`
	ChannelID string `json:"channelID"`
	Amount    int    `json:"amount"`
}

// UserBitsHistory is obsolete
type UserBitsHistory struct {
	Change int       `json:"change"`
	Reason string    `json:"reason"`
	Date   time.Time `json:"date"`
}

// UserBitsWithHistory is obsolete
type UserBitsWithHistory struct {
	UserBits `bson:",inline"`
	History  []UserBitsHistory `json:"history"`
}
