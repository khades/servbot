package models

import (
	"time"
)

type TimerMessageHistory struct {
	TimerMessageInfo `bson:",inline"`
	User             string
	Date             time.Time
}
