package models

import "time"

type UserBits struct {
	User     string `json:"user"`
	UserID   string `json:"userID"`
	ChanneID string `json:"channelID"`
	Amount   int    `json:"amount"`
}

type UserBitsHistory struct {
	Change int       `json:"change"`
	Reason string    `json:"reason"`
	Date   time.Time `json:"date"`
}

type UserBitsWithHistory struct {
	UserBits `bson:",inline"`
	History  []UserBitsHistory `json:"history"`
}
