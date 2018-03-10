package models

//FollowerCursor struct descibes last follower cursor that was processed on specified channel
type FollowerCursor struct {
	ChannelID string
	Cursor    string
}
