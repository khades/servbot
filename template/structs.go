package template

import "time"

// TemplateInfo describes info about chat template WITHOUT history
type TemplateInfo struct {
	ChannelID        string `json:"channelID"`
	CommandName      string `json:"commandName"`
	TemplateInfoBody `bson:",inline"`
}

// TemplateInfoBody describes template settings
type TemplateInfoBody struct {
	AliasTo           string                    `json:"aliasTo"`
	Template          string                    `json:"template"`

}

// TemplateInfoWithHistory describes full template with edit history
type TemplateInfoWithHistory struct {
	TemplateInfo `bson:",inline"`
	History      []TemplateHistory `json:"history"`
}

// TemplateHistory describes edit history of chat command
type TemplateHistory struct {
	TemplateInfoBody `bson:",inline"`
	User             string    `json:"user"`
	UserID           string    `json:"userID"`
	Date             time.Time `json:"date"`
}
