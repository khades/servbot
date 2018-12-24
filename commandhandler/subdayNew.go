package commandhandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"strings"

	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
)

func  (service *CommandHandler) subdayNew(channelInfo *channelInfo.ChannelInfo, chatMessage *chatMessage.ChatMessage, chatCommand models.ChatCommand) {

	if chatMessage.IsMod == false {
		// twitchIRCClient.SendPublic(&models.OutgoingMessage{
		// 	Channel: chatMessage.Channel,
		// 	Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationYoureNotModerator,
		// 	User:    chatMessage.User})
		return
	}
	if channelInfo.SubdayIsActive == true {
		service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationAlreadyExists,
			User:    chatMessage.User})
		return
	}
	subsOnly := true
	subdayName := strings.TrimSpace(strings.Replace(chatCommand.Body, "allowNonSubs=true", "", 1))
	if strings.Contains(chatCommand.Body, "allowNonSubs=true") {
		subsOnly = false
	}
	created, _ := service.subdayService.Create(&chatMessage.ChannelID, subsOnly, &subdayName)
	if created == true {
		service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationSuccess,
			User:    chatMessage.User})
		return
	}
	service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationGeneralError,
		User:    chatMessage.User})
	return
}
