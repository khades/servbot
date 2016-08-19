package models

// ConfigModel defines nasfasf
type ConfigModel struct {
	OauthKey    string `json:",omitempty"`
	BotUserName string `json:",omitempty"`
	Channels    []string
	DbName      string `json:",omitempty"`
}
