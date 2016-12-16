package models

import "time"

type SongRequest struct {
	User             string
	Channel          string
	inQueue          bool
	Paused           bool
	PlayingNow       bool
	VideoID          string
	Length           time.Duration
	EstimatedEndTime time.Time
}
