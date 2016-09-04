package commandHandlers

import (
	"strings"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// Alias sets copy of other command
func Alias(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	if chatMessage.IsMod {
		commandName := ""
		aliasTo := ""
		template := ""
		separator := strings.Index(chatCommand.Body, "=")
		if separator != -1 {
			commandName = chatCommand.Body[:separator]
			aliasTo = strings.TrimSpace(chatCommand.Body[separator+1:])
			result, error := repos.GetChannelTemplate(chatMessage.Channel, aliasTo)
			if error == nil {
				template = result.Template
			}
		} else {
			commandName = chatCommand.Body
			aliasTo = chatCommand.Body
		}
		templateError := Template.updateTemplate(chatMessage.Channel, commandName, template)
		if templateError == nil {
			repos.PutChannelTemplate(chatMessage.User, chatMessage.Channel, commandName, aliasTo, template)
			repos.PushCommandsForChannel(chatMessage.Channel)
			ircClient.SendPublic(models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание алиaса: Ну в принципе готово VoHiYo",
				User:    chatMessage.User})
		} else {
			ircClient.SendPublic(models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    "Создание алиaса: Исходная команда имеет невалидный шаблон etmSad",
				User:    chatMessage.User})
		}
	} else {
		ircClient.SendPublic(models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание алиaса: Вы не модер SMOrc",
			User:    chatMessage.User})
	}
}
