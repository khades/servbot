package models

// OutgoingDebouncedMessage describes info about chat command duh
type OutgoingDebouncedMessage struct {
	*OutgoingMessage
	Command string
}
