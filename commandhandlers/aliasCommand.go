package commandhandlers

import (
	"strings"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// aliasCommand creates alias of existing command by copying its content and setting that it is alias to a parent command
func aliasCommand(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	if chatMessage.IsMod {
		commandName := ""
		aliasTo := ""
		separator := strings.Index(chatCommand.Body, "=")
		if separator != -1 {
			commandName = chatCommand.Body[:separator]
			aliasTo = strings.TrimSpace(chatCommand.Body[separator+1:])

		} else {
			commandName = chatCommand.Body
			aliasTo = chatCommand.Body
		}
		repos.SetChannelTemplateAlias(&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &aliasTo)

		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание алиaса: Ну в принципе готово VoHiYo",
			User:    chatMessage.User})

	} 
	// else {
	// 	ircClient.SendPublic(&models.OutgoingMessage{
	// 		Channel: chatMessage.Channel,
	// 		Body:    "Создание алиaса: Вы не модер SMOrc",
	// 		User:    chatMessage.User})
	// }
	
}
