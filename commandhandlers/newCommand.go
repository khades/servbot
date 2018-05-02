package commandhandlers

import (
	"strings"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// newCommand creates or modifies mustache templates for channel
func newCommand(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
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
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CommandCyclingIsForbidden,
				User:    chatMessage.User})
			return
		}
		if commandName == "" {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).EmptyCommandNameIsForbidden,
				User:    chatMessage.User})
			return
		}
		if commandName == "new" || commandName == "alias" {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).ReservedCommandNameIsForbidden,
				User:    chatMessage.User})
			return
		}
		templateBody := models.TemplateInfoBody{
			Template: template}
		templateError := repos.SetChannelTemplate(&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &templateBody)
		if templateError == nil {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CommandCreationSuccess,
				User:    chatMessage.User})
		} else {
			ircClient.SendPublic(&models.OutgoingMessage{
				Channel: chatMessage.Channel,
				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).InvalidCommandTemplate,
				User:    chatMessage.User})
		}

	}
}
