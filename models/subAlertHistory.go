package models

import "time"

// SubAlertHistory describes all manipulations with subalerts
type SubAlertHistory struct {
	User         string    `json:"user"`
	UserID       string    `json:"userID"`
	Date         time.Time `json:"date"`
	Enabled      bool      `json:"enabled"`
	SubMessage   string    `json:"subMessage"`
	ResubMessage string    `json:"resubMessage"`
	RepeatBody   string    `json:"repeatBody"`
}
