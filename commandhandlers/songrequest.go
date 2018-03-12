package commandhandlers

import (
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func songRequestAdd(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	id := &chatCommand.Body
	result := repos.AddSongRequest(&chatMessage.User, chatMessage.IsSub, &chatMessage.UserID, &chatMessage.ChannelID, id)
	ircClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    result.Error(),
		User:    chatMessage.User})
	return
}