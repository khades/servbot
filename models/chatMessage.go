package models

import (
	"strings"
	"time"
)

// ChatMessage describes processed twitch message with essential information on it
type ChatMessage struct {
	Channel          string
	User             string
	IsMod            bool
	IsSub            bool
	IsPrime          bool
	Date             time.Time
	MessageBody      string            `bson:",omitempty"`
	SubscriptionInfo *SubscriptionInfo `bson:",omitempty"`
	BanInfo          *BanInfo          `bson:",omitempty"`
}

// GetCommand method checks if message starts from ! and returns body of command if it is command
func (chatMessage ChatMessage) GetCommand() (ChatCommand, bool) {
	chatCommand := ChatCommand{}
	isCommand := strings.HasPrefix(chatMessage.MessageBody, "!")
	if isCommand {
		spaceIndex := strings.Index(chatMessage.MessageBody, " ")
		if spaceIndex != -1 {
			chatCommand = ChatCommand{
				Command: chatMessage.MessageBody[1:spaceIndex],
				Body:    chatMessage.MessageBody[spaceIndex+1:]}
		} else {
			chatCommand = ChatCommand{
				Command: chatMessage.MessageBody[1:]}
		}
	}
	return chatCommand, isCommand
}
