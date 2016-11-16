package models

import "time"

// TimedMessageInfo describes
type TimerMessageInfo struct {
	ID                string
	Channel           string
	Period            int
	MessageThreshhold int
	Body              string
	LastRun           time.Time
}
