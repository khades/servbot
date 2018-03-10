package models

import "time"


// AutoMessageHistory struct describes edit history of automessage
type AutoMessageHistory struct {
	User          string        `valid:"required" json:"user"`
	UserID        string        `valid:"required" json:"userID"`
	Game          string        `json:"game"`
	Date          time.Time     `json:"date"`
	Message       string        `json:"message"`
	MessageLimit  int           `json:"messageLimit"`
	DurationLimit time.Duration `json:"durationLimit"`
}
