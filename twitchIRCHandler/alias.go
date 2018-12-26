package twitchIRCHandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/twitchIRC"
	"strings"
)


// aliasCommand creates alias of existing command by copying its content and setting that it is alias to a parent command
func (service *TwitchIRCHandler) alias(
	channelInfo *channelInfo.ChannelInfo,
	chatMessage *chatMessage.ChatMessage,
	chatCommand chatMessage.ChatCommand,
	twitchIRCClient *twitchIRC.Client) {
	if chatMessage.IsMod {
		commandName := ""
		aliasTo := ""
		separator := strings.Index(chatCommand.Body, "=")
		if separator != -1 {
			commandName = strings.ToLower(strings.Join(strings.Fields(chatCommand.Body[:separator]), ""))
			aliasTo = strings.ToLower(strings.Join(strings.Fields(chatCommand.Body[separator+1:]), ""))
		} else {
			commandName =  strings.ToLower(strings.Join(strings.Fields(chatCommand.Body), ""))
			aliasTo =  strings.ToLower(strings.Join(strings.Fields(chatCommand.Body), ""))
		}
		service.templateService.SetAlias(&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &aliasTo)

		twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).AliasCreationSuccess,
			User:    chatMessage.User})

	}
	// else {
	// 	twitchIRC.SendPublic(&models.OutgoingMessage{
	// 		Channel: chatMessage.Channel,
	// 		Body:    "Создание алиaса: Вы не модер SMOrc",
	// 		User:    chatMessage.User})
	// }

}

