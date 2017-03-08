package models

// TemplateInfo describes info about chat template WITHOUT history
type TemplateInfo struct {
	CommandName string `json:"commandName"`
	ChannelID   string `json:"channelID"`
	AliasTo     string `json:"aliasTo"`
	Template    string `json:"template"`
}
