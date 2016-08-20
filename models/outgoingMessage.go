package models

// OutgoingMessage describes info about chat command duh
type OutgoingMessage struct {
	Channel string
	Body    string
	User    string
}
