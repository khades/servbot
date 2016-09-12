package models

// OutgoingMessage defines message that is send to a server
type OutgoingMessage struct {
	Channel string
	Body    string
	User    string
}
