package models

import "time"

type StreamStatusGameHistory struct {
	Game  string    `bson:",omitempty"`
	Start time.Time `bson:",omitempty"`
}
type GamesHistory []StreamStatusGameHistory

func (history GamesHistory) Len() int {
	return len(history)
}
func (history GamesHistory) Less(i, j int) bool {
	return history[i].Start.Before(history[j].Start)
}

func (history GamesHistory) Swap(i, j int) {
	history[i], history[j] = history[j], history[i]
}

// StreamStatus Describes info about stream, when started, what game and title is, and if it is online
type StreamStatus struct {
	Online         bool
	Game           string       `bson:",omitempty"`
	Title          string       `bson:",omitempty"`
	Start          time.Time    `bson:",omitempty"`
	LastOnlineTime time.Time    `bson:",omitempty"`
	Viewers        int          `bson:",omitempty"`
	GamesHistory   GamesHistory `bson:",omitempty"`
}
