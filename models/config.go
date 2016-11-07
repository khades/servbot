package models

// Config defines nasfasf
type Config struct {
	OauthKey    string `valid:"required"`
	BotUserName string `valid:"required"`
	Channels    []string
	DbName      string `valid:"required"`
	ClientKey   string `valid:"required"`
}
