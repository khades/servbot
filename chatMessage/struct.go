package chatMessage

import (
	"strings"
	"time"
)

// ChatMessage describes processed twitch message with essential information on it
type ChatMessage struct {
	MessageStruct `bson:",inline"`
	Channel       string
	ChannelID     string
	User          string
	UserID        string
	IsMod         bool
	IsSub         bool
	IsPrime       bool
}

type MessageStruct struct {
	Date        time.Time `json:"date"`
	Username    string    `json:"username"`
	MessageBody string    `json:"messageBody"`
	MessageType string    `json:"messageType"`
	BanLength   int       `json:"banLength"`
	BanReason   string    `json:"banReason"`
	BanIssuer   string    `json:"banIssuer"`
	BanIssuerID string    `json:"banIssuerID"`
}

// ChatCommand describes info about incoming chat command
type ChatCommand struct {
	Command string
	Body    string
}

// GetCommand method checks if message starts from ! and returns body of command if it is command
func (chatMessage ChatMessage) GetCommand() (ChatCommand, bool) {
	chatCommand := ChatCommand{}
	isCommand := strings.HasPrefix(chatMessage.MessageBody, "!")
	if isCommand {
		spaceIndex := strings.Index(chatMessage.MessageBody, " ")
		if spaceIndex != -1 {
			chatCommand = ChatCommand{
				Command: strings.ToLower(strings.TrimSpace(chatMessage.MessageBody[1:spaceIndex])),
				Body:    chatMessage.MessageBody[spaceIndex+1:]}
		} else {
			chatCommand = ChatCommand{
				Command: strings.ToLower(strings.TrimSpace(chatMessage.MessageBody[1:]))}
		}
	}
	return chatCommand, isCommand
}
