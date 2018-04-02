package models

import "time"

//FollowerCursor struct descibes last follower cursor that was processed on specified channel
type FollowerCursor struct {
	ChannelID string
	Cursor    time.Time
}
