package models

import "time"

// TemplateHistory describes history of chat command
type TemplateHistory struct {
	TemplateInfo `bson:",inline"`
	User         string
	Date         time.Time
}
