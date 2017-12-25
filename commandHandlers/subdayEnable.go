package commandHandlers

import (
	//"strings"
	//"time"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func SubdayEnable(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {

	if chatMessage.IsMod == false {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание сабдня: Вы не модер SMOrc",
			User:    chatMessage.User})
		return
	}
	repos.SetSubdayEnabled(&chatMessage.ChannelID, true)
	ircClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    "Сабдеи включены",
		User:    chatMessage.User})
	return
}