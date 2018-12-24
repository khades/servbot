package gameResolve

import "time"

// Game describes relation of game and gameID on twitch
type Game struct {
	Game   string    `json:"game"`
	GameID string    `json:"gameID"`
	Date   time.Time `json:"date"`
}
