package models

import (
	"time"
)

// TimerMessageHistory is history-item for TimerMessageInfo
type TimerMessageHistory struct {
	TimerMessageInfo `bson:",inline"`
	User             string
	Date             time.Time
}
