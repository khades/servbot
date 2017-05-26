package models

import "time"

// TemplateHistory describes history of chat command
type TemplateHistory struct {
	TemplateInfoBody `bson:",inline"`
	User             string    `json:"user"`
	UserID           string    `json:"userID"`
	Date             time.Time `json:"date"`
}
