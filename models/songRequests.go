package models

import (
	"time"
)
// SongRequest struct descibes one song request
type SongRequest struct {
	User    string        `json:"user"`
	UserID  string        `json:"userID"`
	Date    time.Time     `json:"date"`
	VideoID string        `json:"videoID"`
	Length  time.Duration `json:"length"`
	Title   string        `json:"title"`
	Order   int32         `json:"order"`
}

// ChannelSongRequestSettings struct descibes current settings for songrequest on channel
type ChannelSongRequestSettings struct {
	OnlySubs           bool  `json:"onlySubs"`
	PlaylistLength     int   `json:"playlistLength"`
	MaxVideoLength     int   `json:"maxVideoLength"`
	MaxRequestsPerUser int   `json:"maxRequestsPerUser"`
	VideoViewLimit     int64 `json:"videoViewLimit"`
}

// ChannelSongRequest describes all info about song request for channel
type ChannelSongRequest struct {
	ChannelID string                     `json:"channelID"`
	Settings  ChannelSongRequestSettings `json:"settings"`
	Requests  []SongRequest              `json:"requests"`
}
