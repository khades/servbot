package models

type UserBits struct {
	User     string `json:"user"`
	UserID   string `json:"userID"`
	ChanneID string `json:"channelID"`
	Amount   int    `json:"amount"`
}
