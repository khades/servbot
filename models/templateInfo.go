package models

import "time"

// TemplateInfo describes info about chat template WITHOUT history
type TemplateInfo struct {
	ChannelID        string `json:"channelID"`
	CommandName      string `json:"commandName"`
	TemplateInfoBody `bson:",inline"`
}

// TemplateInfoBody describes template settings 
type TemplateInfoBody struct {
	ShowOnline        bool                      `json:"showOnline"`
	ShowOffline       bool                      `json:"showOffline"`
	PreventDebounce   bool                      `json:"preventDebounce"`
	PreventRedirect   bool                      `json:"preventRedirect"`
	OnlyPrivate       bool                      `json:"onlyPrivate"`
	AliasTo           string                    `json:"aliasTo"`
	Template          string                    `json:"template"`
	IntegerRandomizer TemplateIntegerRandomizer `json:"integerRandomizer"`
	StringRandomizer  TemplateStringRandomizer  `json:"stringRandomizer"`
}
// TemplateIntegerRandomizer describes template integer randomizator settings
type TemplateIntegerRandomizer struct {
	Enabled      bool `json:"enabled"`
	TimeoutAfter bool `json:"timeoutAfter"`
	LowerLimit   int  `json:"lowerLimit"`
	UpperLimit   int  `json:"upperLimit"`
}

// TemplateStringRandomizer describes template string randomizator settings
type TemplateStringRandomizer struct {
	Enabled bool     `json:"enabled"`
	Strings []string `json:"strings"`
}

// TemplateExtendedObject describes extended channelInfo object for template processing
type TemplateExtendedObject struct {
	ChannelInfo
	RandomInteger          int
	RandomIntegerIsMinimal bool
	RandomIntegerIsMaximal bool
	RandomIngegerIsZero    bool
	RandomString           string
	IsMod                  bool
	IsSub                  bool
	CommandBody            string
	CommandBodyIsEmpty     bool
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
