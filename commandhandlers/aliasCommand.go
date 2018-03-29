package commandhandlers

import (
	"strings"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/l10n"
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
			commandName = strings.ToLower(chatCommand.Body[:separator])
			aliasTo = strings.ToLower(strings.TrimSpace(chatCommand.Body[separator+1:]))

		} else {
			commandName = chatCommand.Body
			aliasTo = chatCommand.Body
		}
		repos.SetChannelTemplateAlias(&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &aliasTo)

		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).AliasCreationSuccess,
			User:    chatMessage.User})

	}
	// else {
	// 	ircClient.SendPublic(&models.OutgoingMessage{
	// 		Channel: chatMessage.Channel,
	// 		Body:    "Создание алиaса: Вы не модер SMOrc",
	// 		User:    chatMessage.User})
	// }

}
