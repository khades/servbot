package commandHandlers

import (
	//"strings"
	//"time"
	"log"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func SongRequestAdd(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	id := &chatCommand.Body
	log.Println(chatCommand.Body)
	result := repos.AddSongRequest(&chatMessage.User, chatMessage.IsSub, &chatMessage.UserID, &chatMessage.ChannelID, id)
	repos.SetSubdayEnabled(&chatMessage.ChannelID, false)
	ircClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    result.Error(),
		User:    chatMessage.User})
	return
}