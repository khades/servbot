package models

// TemplateInfo describes info about chat template WITHOUT history
type TemplateInfo struct {
	ChannelID        string `json:"channelID"`
	CommandName      string `json:"commandName"`
	TemplateInfoBody `bson:",inline"`
}

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
type TemplateIntegerRandomizer struct {
	Enabled      bool `json:"enabled"`
	TimeoutAfter bool `json:"timeoutAfter"`
	LowerLimit   int  `json:"lowerLimit"`
	UpperLimit   int  `json:"upperLimit"`
}

type TemplateStringRandomizer struct {
	Enabled bool     `json:"enabled"`
	Strings []string `json:"strings"`
}
