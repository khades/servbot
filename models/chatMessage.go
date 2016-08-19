package models

import (
	"strings"
	"time"
)

// ChatMessage describes processed twitch message with essential information on it
type ChatMessage struct {
	Channel     string
	Username    string
	MessageBody string
	IsMod       bool
	IsSub       bool
	Date        time.Time
}

func (chatMessage ChatMessage) isCommand() (bool, ChatCommand) {
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
	return isCommand, chatCommand
}
