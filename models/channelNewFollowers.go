package models

//ChannelNewFollowers describes new followers on channel
type ChannelNewFollowers struct {
	Channel string
	ChannelID string
	Followers []string
}