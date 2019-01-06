package event

import "time"

type EventType = string

const SUB EventType = "SUB"
const RESUB EventType = "RESUB"
const FOLLOW EventType = "FOLLOW"
const BITS EventType = "BITS"
const DONATION EventType = "DONATION"

type Events struct {
	ChannelID  string
	ViewedTill time.Time
	Events     []Event
}

type Event struct {
	User     string    `json:"user"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date"`
	Type     EventType `json:"type"`
	Amount   int       `json:"amount"`
	Currency string    `json:"currency"`
}
