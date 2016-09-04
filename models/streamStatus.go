package models

import "time"

// StreamStatus Describes info about stream, when started, what game and title is, and if it is online
type StreamStatus struct {
	Online      bool
	Description string    `bson:",omitempty"`
	Game        string    `bson:",omitempty"`
	Title       string    `bson:",omitempty"`
	Start       time.Time `bson:",omitempty"`
}
