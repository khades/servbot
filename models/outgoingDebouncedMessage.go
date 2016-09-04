package models

// OutgoingDebouncedMessage describes info about chat command duh
type OutgoingDebouncedMessage struct {
	Message    OutgoingMessage
	Command    string
	Redirected bool
}
