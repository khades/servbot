package twitchIRCHandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/twitchIRC"
	"strings"
)

func  (service *TwitchIRCHandler) subdayNew(
	channelInfo *channelInfo.ChannelInfo,
	chatMessage *chatMessage.ChatMessage,
	chatCommand chatMessage.ChatCommand,
	twitchIRCClient *twitchIRC.Client) {

	if chatMessage.IsMod == false {
		// twitchIRC.SendPublic(&models.OutgoingMessage{
		// 	Channel: chatMessage.Channel,
		// 	Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationYoureNotModerator,
		// 	User:    chatMessage.User})
		return
	}
	if channelInfo.SubdayIsActive == true {
		twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
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
	created, _ := service.subdayService.Create(channelInfo, subsOnly, &subdayName)
	if created == true {
		twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationSuccess,
			User:    chatMessage.User})
		return
	}
	twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationGeneralError,
		User:    chatMessage.User})
	return
}
