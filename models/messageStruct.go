package models

import "time"

type MessageStruct struct {
	Date        time.Time `json:"date"`
	Username    string    `json:"username"`
	MessageBody string    `json:"messageBody"`
	MessageType string    `json:"messageType"`
	BanLength   int       `json:"banLength"`
	BanReason   string    `json:"banReason"`
}
