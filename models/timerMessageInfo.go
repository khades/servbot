package models

import "time"

// TimerMessageInfo describes moobot-like automatic messages to chat
type TimerMessageInfo struct {
	ID                string
	Channel           string
	Period            int
	MessageThreshhold int
	Body              string
	LastRun           time.Time
}
