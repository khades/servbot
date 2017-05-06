package commandHandlers

import (
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
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Запрещено зацикливать команды",
				User:    chatMessage.User})
			return
		}
		if commandName == "" {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Запрещено создавать пустые команды",
				User:    chatMessage.User})
			return
		}
		if commandName == "new" {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Запрещено создавать команды для зарезервированных слов",
				User:    chatMessage.User})
			return
		}
		templateError := repos.SetChannelTemplate(&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &template)
		if templateError == nil {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Ну в принципе готово VoHiYo",
				User:    chatMessage.User})
		} else {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание команды: Невалидный шаблон для команды etmSad",
				User:    chatMessage.User})
		}

	} else {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание алиaса: Вы не модер SMOrc",
			User:    chatMessage.User})
	}
}
