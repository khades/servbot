package models

import "time"

// SubAlertHistory describes all manipulations with subalerts
type SubAlertHistory struct {
	SubAlertInfo `bson:",inline"`
	User         string
	Date         time.Time
}
