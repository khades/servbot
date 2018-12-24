package twitchIRCClient

import (
	irc "gopkg.in/irc.v2"
)

// type TwitchIRCHandle
type TwitchIRCHandle func(client *TwitchIRCClient, message *irc.Message)


// HandlerFunc is a simple wrapper around a function which allows it
// to be used as a Handler.
//type HandlerFunc func(*Client, *Message)
