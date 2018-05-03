package commandhandlers

import (
	"log"
	"strings"
	"time"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func subdayNew(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	log.Println("CREATING SUBDAY")

	if chatMessage.IsMod == false {
		// ircClient.SendPublic(&models.OutgoingMessage{
		// 	Channel: chatMessage.Channel,
		// 	Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationYoureNotModerator,
		// 	User:    chatMessage.User})
		return
	}
	if channelInfo.SubdayIsActive == true {
		ircClient.SendPublic(&models.OutgoingMessage{
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
	if subdayName == "" {
		subdayName = l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationPrefix + time.Now().Format(time.UnixDate)
	}
	created, _ := repos.CreateNewSubday(&chatMessage.ChannelID, subsOnly, &subdayName)
	if created == true {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationSuccess,
			User:    chatMessage.User})
		return
	}
	ircClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayCreationGeneralError,
		User:    chatMessage.User})
	return
}
