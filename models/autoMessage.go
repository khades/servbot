package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type AutoMessage struct {
	ID                bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ChannelID         string        `json:"channel"`
	Message           string        `json:"message"`
	MessageLimit      int           `json:"messageLimit"`
	DurationLimit     time.Duration `json:"durationLimit"`
	MessageThreshold  int           `json:"messageThreshold"`
	DurationThreshold time.Time     `json:"durationThreshold"`
}

// func (autoMessage AutoMessage) ToOutgoingMessage() *OutgoingMessage {
// 	return &OutgoingMessage{Channel: autoMessage.Channel, Body: autoMessage.Message}

// }
