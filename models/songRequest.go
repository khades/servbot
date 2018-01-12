package models

import (
	"time"

)

type SongRequest struct {
	User      string
	UserID    string
	Date      time.Time
	VideoID string
	Length  time.Duration
	Title string
}

type ChannelSongRequestSettings struct {
	OnlySubs bool
	PlaylistLength int
	MaxVideoLength int
	MaxRequestsPerUser int
}
type ChannelSongRequest struct {
	ChannelID string
	Settings ChannelSongRequestSettings
	Requests []SongRequest
}