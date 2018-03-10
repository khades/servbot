package models

// ChannelUser struct describes user on one channel 
type ChannelUser struct {
	ChannelID      string   `json:"channelID"`
	User           string   `json:"user"`
	UserID         string   `json:"userID"`
	KnownNicknames []string `json:"knownNicknames"`
}
