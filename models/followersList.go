package models

//FollowersList struct describes all seen followers on specified channel
type FollowersList struct {
	ChannelID string
	Followers []string
}
