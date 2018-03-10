package commandhandlers

import (
	"log"
	"strings"
	"time"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func subdayNew(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	
	
	if chatMessage.IsMod == false {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Создание сабдня: Вы не модер SMOrc",
			User:    chatMessage.User})
		return
	}
	_, subdayError := repos.GetLastActiveSubday(&chatMessage.ChannelID)
	if subdayError == nil {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Сабдей уже существует",
			User:    chatMessage.User})
		return
	}
	subsOnly := true
	subdayName := strings.TrimSpace(strings.Replace(chatCommand.Body, "allowNonSubs=true", "", 1))
	if strings.Contains(chatCommand.Body, "allowNonSubs=true") {
		subsOnly = false
	}
	if subdayName == "" {
		subdayName = "Сабдей, созданный " + time.Now().Format(time.UnixDate)
	}
	created := repos.CreateNewSubday(&chatMessage.ChannelID, subsOnly, &subdayName)
	if created == true {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    "Сабдей создан! VoHiYo",
			User:    chatMessage.User})
		return

	}
	ircClient.SendPublic(&models.OutgoingMessage{
		Channel: chatMessage.Channel,
		Body:    "Что-то пошло не так etmSad",
		User:    chatMessage.User})
	return
}
