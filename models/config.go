package models

// Config defines nasfasf
type Config struct {
	OauthKey     string `valid:"required"`
	BotUserName  string `valid:"required"`
	Channels     []string
	DbName       string `valid:"required"`
	ClientID     string `valid:"required"`
	ClientSecret string `valid:"required"`
	AppOauthURL  string `valid:"required"`
	AppURL       string `valid:"required"`
	Debug        bool   `valid:"required"`
}
