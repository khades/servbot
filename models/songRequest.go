package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SongRequest struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	User    string
	Channel string
	Date    time.Time
	// inQueue    bool
	// PlayingNow bool
	VideoID string
	// Length           time.Duration
	// EstimatedEndTime time.Time
}
