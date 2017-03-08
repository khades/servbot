package models

import "time"

type AutoMessageHistory struct {
	User   string `valid:"required" json:"user"`
	UserID string `valid:"required" json:"userID"`

	Date          time.Time     `json:"date"`
	Message       string        `json:"message"`
	MessageLimit  int           `json:"messageLimit"`
	DurationLimit time.Duration `json:"durationLimit"`
}
