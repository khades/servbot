package models

// OutgoingDebouncedMessage defines message that is send to a server with debounce abilities
type OutgoingDebouncedMessage struct {
	Message    OutgoingMessage
	Command    string
	RedirectTo string
}
