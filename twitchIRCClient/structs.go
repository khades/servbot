package twitchIRCClient

type OutgoingMessage struct {
	Channel string
	Body    string
	User    string
}
// OutgoingDebouncedMessage defines message that is send to a server with debounce abilities
type OutgoingDebouncedMessage struct {
	Message    OutgoingMessage
	Command    string
	RedirectTo string
}
