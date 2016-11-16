package models

import "time"

type SubAlertHistory struct {
	SubAlertInfo `bson:",inline"`
	User         string
	Date         time.Time
}
