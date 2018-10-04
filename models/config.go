package models

// Config struct describes configuration of that chatbot instance
type Config struct {
	APIKey       string `valid:"required"`
	OauthKey     string `valid:"required"`
	BotUserName  string `valid:"required"`
	BotUserID    string `valid:"required"`
	Channels     []string
	ChannelIDs   []string
	ClientID     string `valid:"required"`
	ClientSecret string `valid:"required"`
	AppOauthURL  string `valid:"required"`
	AppURL       string `valid:"required"`
	Debug        bool
	VkClientKey  string
	YoutubeKey   string
}
