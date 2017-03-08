package models

import "time"

// TemplateHistory describes history of chat command
type TemplateHistory struct {
	AliasTo  string    `json:"aliasTo"`
	Template string    `json:"template"`
	User     string    `json:"user"`
	UserID   string    `json:"userID"`
	Date     time.Time `json:"date"`
}
