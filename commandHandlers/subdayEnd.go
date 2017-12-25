package commandHandlers

import (
	//"strings"
	//"time"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func SubdayEnd(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	channelInfo, error := repos.GetChannelInfo(&chatMessage.ChannelID)
	if error != nil || channelInfo.SubdayEnabled == false {
		return
	}
	
	if chatMessage.IsMod == false {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание сабдня: Вы не модер SMOrc",
			User:    chatMessage.User})
		return
	}
	_, subdayError := repos.GetLastActiveSubday(&chatMessage.ChannelID)
	if subdayError == nil {
		repos.CloseActiveSubday(&chatMessage.ChannelID)
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Сабдей закрыт",
			User:    chatMessage.User})
		return
	}
	ircClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    "Открытого сабдея не существует",
		User:    chatMessage.User})
	return
	
}
