package followers

import "time"

//FollowersList struct describes all seen followers on specified channel
type FollowersList struct {
	ChannelID string
	Followers []Follower
}

//Follower struct describes one follower and his follow date
type Follower struct {
	ChannelID  string
	UserID     string
	Date       time.Time
	IsFollower bool
}
