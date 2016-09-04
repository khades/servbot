package commandHandlers

import (
	"strings"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// New works with new command that adds template
func New(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	if chatMessage.IsMod {
		if strings.HasPrefix(chatCommand.Body, "!") || strings.HasPrefix(chatCommand.Body, ".") || strings.HasPrefix(chatCommand.Body, "/") {
			ircClient.SendPublic(models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Запрещено зацикливать команды",
				User:    chatMessage.User})
		} else {
			commandName := ""
			template := ""
			separator := strings.Index(chatCommand.Body, "=")
			if separator != -1 {
				commandName = chatCommand.Body[:separator]
				template = chatCommand.Body[separator+1:]
			} else {
				commandName = chatCommand.Body
			}
			repos.PutChannelTemplate(chatMessage.User, chatMessage.Channel, commandName, template)
			Template.update(chatMessage.Channel, commandName, template)
			ircClient.SendPublic(models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Ну в принципе готово VoHiYo",
				User:    chatMessage.User})
		}

	} else {
		ircClient.SendPublic(models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Вы не модер SMOrc",
			User:    chatMessage.User})
	}
}
