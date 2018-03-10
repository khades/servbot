package models

//ChannelWithID struct describes channel to channelID relation
type ChannelWithID struct {
	Channel   string `json:"channel"`
	ChannelID string `json:"channelID"`
}
