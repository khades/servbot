package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// AutoMessage struct describes one automessage item on specified channel
type AutoMessage struct {
	ID                bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ChannelID         string        `json:"channelID"`
	Message           string        `json:"message"`
	MessageLimit      int           `json:"messageLimit"`
	DurationLimit     time.Duration `json:"durationLimit"`
	MessageThreshold  int           `json:"messageThreshold"`
	DurationThreshold time.Time     `json:"durationThreshold"`
	Game              string        `json:"game"`
}
