package models

type ChannelUsers struct {
	ChannelID      string   `json:"channelID"`
	User           string   `json:"user"`
	UserID         string   `json:"userID"`
	KnownNicknames []string `json:"knownNicknames"`
}
