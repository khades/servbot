package commandhandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/template"
	"strings"

	"github.com/khades/servbot/twitchIRCClient"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
)

// newCommand creates or modifies mustache templates for channel
func  (service *CommandHandler) newCommand(channelInfo *channelInfo.ChannelInfo,  chatMessage *chatMessage.ChatMessage, chatCommand models.ChatCommand) {
	if chatMessage.IsMod {
		commandName := ""
		template := ""
		separator := strings.Index(chatCommand.Body, "=")
		if separator != -1 {

			commandName = strings.ToLower(strings.Join(strings.Fields(chatCommand.Body[:separator]), ""))
			template = strings.TrimSpace(chatCommand.Body[separator+1:])
		} else {
			commandName = strings.ToLower(strings.Join(strings.Fields(chatCommand.Body), ""))
		}
		if strings.HasPrefix(template, "!") || strings.HasPrefix(template, ".") || strings.HasPrefix(template, "/") {
			service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CommandCyclingIsForbidden,
				User:    chatMessage.User})
			return
		}
		if commandName == "" {
			service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).EmptyCommandNameIsForbidden,
				User:    chatMessage.User})
			return
		}
		if commandName == "new" || commandName == "alias" {
			service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).ReservedCommandNameIsForbidden,
				User:    chatMessage.User})
			return
		}
		templateBody := template.TemplateInfoBody{
			Template: template}
		templateError := service.templateService.Set(&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &templateBody)
		if templateError == nil {
			service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CommandCreationSuccess,
				User:    chatMessage.User})
		} else {
			service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).InvalidCommandTemplate,
				User:    chatMessage.User})
		}

	}
}
