package models

import "time"

// TemplateHistory describes history of chat command
type TemplateHistory struct {
	AliasTo  string
	Template string
	User     string
	Date     time.Time
}
