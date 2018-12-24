package webhookAPI

import "time"

type twitchPubSubFollows struct {
	Data      twitchPubSubFollower `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}
type twitchPubSubFollower struct {
	ChannelID string `json:"to_id"`
	UserID    string `json:"from_id"`
}

type twitchPubSubStreams struct {
	Data []twitchPubSubStream `json:"data"`
}
type twitchPubSubStream struct {
	ChannelID string    `json:"user_id"`
	GameID    string    `json:"game_id"`
	Title     string    `json:"title"`
	Viewers   int       `json:"viewer_count"`
	Start     time.Time `json:"started_at"`
}
