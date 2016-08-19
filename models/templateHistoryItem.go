package models

import "time"

// TemplateHistoryItem describestem of chat template history
type TemplateHistoryItem struct {
	Template string
	User     string
	Date     time.Time
}
