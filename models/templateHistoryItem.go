package models

import "time"

// TemplateHistoryItem describes history of chat command
type TemplateHistoryItem struct {
	Template string
	User     string
	Date     time.Time
}
