package twitchIRCHandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/template"
	"github.com/khades/servbot/twitchIRC"
	"strings"
)

func(service *TwitchIRCHandler) new(channelInfo *channelInfo.ChannelInfo,  chatMessage *chatMessage.ChatMessage, chatCommand chatMessage.ChatCommand, twitchIRCClient *twitchIRC.Client) {
	if chatMessage.IsMod {
		commandName := ""
		templateBody := ""
		separator := strings.Index(chatCommand.Body, "=")
		if separator != -1 {

			commandName = strings.ToLower(strings.Join(strings.Fields(chatCommand.Body[:separator]), ""))
			templateBody = strings.TrimSpace(chatCommand.Body[separator+1:])
		} else {
			commandName = strings.ToLower(strings.Join(strings.Fields(chatCommand.Body), ""))
		}
		if strings.HasPrefix(templateBody, "!") || strings.HasPrefix(templateBody, ".") || strings.HasPrefix(templateBody, "/") {
			twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CommandCyclingIsForbidden,
				User:    chatMessage.User})
			return
		}
		if commandName == "" {
			twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).EmptyCommandNameIsForbidden,
				User:    chatMessage.User})
			return
		}
		if commandName == "new" || commandName == "alias" {
			twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).ReservedCommandNameIsForbidden,
				User:    chatMessage.User})
			return
		}
		templateBodyStruct := template.TemplateInfoBody{
			Template: templateBody}
		templateError := service.templateService.Set(&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &templateBodyStruct)
		if templateError == nil {
			twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CommandCreationSuccess,
				User:    chatMessage.User})
		} else {
			twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).InvalidCommandTemplate,
				User:    chatMessage.User})
		}

	}
}

