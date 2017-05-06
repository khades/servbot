package models

import "time"

// SubAlertHistory describes all manipulations with subalerts
type SubAlertHistory struct {
	User     string    `json:"user"`
	UserID   string    `json:"userID"`
	Date     time.Time `json:"date"`
	SubAlert `bson:",inline"`
}
