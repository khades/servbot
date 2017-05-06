package models

import "time"

type StreamStatusGameHistory struct {
	Game  string    `bson:",omitempty"`
	Start time.Time `bson:",omitempty"`
}

// StreamStatus Describes info about stream, when started, what game and title is, and if it is online
type StreamStatus struct {
	Online         bool
	Game           string                    `bson:",omitempty"`
	Title          string                    `bson:",omitempty"`
	Start          time.Time                 `bson:",omitempty"`
	LastOnlineTime time.Time                 `bson:",omitempty"`
	Viewers        int                       `bson:",omitempty"`
	GamesHistory   []StreamStatusGameHistory `bson:",omitempty"`
}
