package autoMessage

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// AutoMessage struct describes one autoMessage item on specified channel
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

// AutoMessageHistory struct describes edit history of autoMessage
type AutoMessageHistory struct {
	User          string        `valid:"required" json:"user"`
	UserID        string        `valid:"required" json:"userID"`
	Game          string        `json:"game"`
	Date          time.Time     `json:"date"`
	Message       string        `json:"message"`
	MessageLimit  int           `json:"messageLimit"`
	DurationLimit time.Duration `json:"durationLimit"`
}

// AutoMessageUpdate struct describes update to autoMessage
type AutoMessageUpdate struct {
	ID            string
	User          string `valid:"required"`
	UserID        string `valid:"required"`
	ChannelID     string `valid:"required"`
	Game          string `json:"game"`
	Message       string `valid:"required" json:"message"`
	MessageLimit  int    `valid:"required" json:"messageLimit"`
	DurationLimit int    `valid:"required" json:"durationlimit"`
}

type AutoMessageWithHistory struct {
	AutoMessage `bson:",inline"`
	History     []AutoMessageHistory `json:"history"`
}
