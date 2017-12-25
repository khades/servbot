package commandHandlers

import (
	//"strings"
	//"time"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func SubdayDisable(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {

	if chatMessage.IsMod == false {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание сабдня: Вы не модер SMOrc",
			User:    chatMessage.User})
		return
	}
	repos.SetSubdayEnabled(&chatMessage.ChannelID, false)
	ircClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    "Сабдеи отключены",
		User:    chatMessage.User})
	return
}