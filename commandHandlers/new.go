package commandHandlers

import (
	"log"
	"strings"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// New creates or modifies mustache templates for channel
func New(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	if chatMessage.IsMod {
		commandName := ""
		template := ""
		separator := strings.Index(chatCommand.Body, "=")
		if separator != -1 {
			commandName = chatCommand.Body[:separator]
			template = strings.TrimSpace(chatCommand.Body[separator+1:])
		} else {
			commandName = chatCommand.Body
		}
		if strings.HasPrefix(template, "!") || strings.HasPrefix(template, ".") || strings.HasPrefix(template, "/") {
			ircClient.SendPublic(models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Запрещено зацикливать команды",
				User:    chatMessage.User})
		} else {
			templateError := Template.updateTemplate(chatMessage.Channel, commandName, template)
			if templateError == nil {
				repos.PutChannelTemplate(chatMessage.User, chatMessage.Channel, commandName, commandName, template)
				Template.updateAliases(chatMessage.Channel, commandName, template)
				repos.PushCommandsForChannel(chatMessage.Channel)
				ircClient.SendPublic(models.OutgoingMessage{
					Channel: chatMessage.Channel,
					Body:    "Создание команды: Ну в принципе готово VoHiYo",
					User:    chatMessage.User})
			} else {
				log.Println(templateError)
				ircClient.SendPublic(models.OutgoingMessage{
					Channel: chatMessage.Channel,
					Body:    "Создание команды: Невалидный шаблон для команды etmSad",
					User:    chatMessage.User})
			}

		}
	} else {
		ircClient.SendPublic(models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание алиaса: Вы не модер SMOrc",
			User:    chatMessage.User})
	}
}
