package models

import "time"

// TwitchStreamStatus describes stream status returned by twitch api
type TwitchStreamStatus struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
	CreatedAt time.Time `json:"started_at"`
	Title string `json:"title"`
	ViewerCount   int `json:"viewer_count"`
}

