package commandHandlers

import (
	"unicode/utf8"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func Vote(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	if utf8.RuneCountInString(chatMessage.MessageBody) < 2 {
		return
	}
	game := chatMessage.MessageBody[1:]
	subday, subdayError := repos.GetLastActiveSubday(&chatMessage.ChannelID)

	if subdayError != nil {
		ircClient.SendDebounced(models.OutgoingDebouncedMessage{
			Message: models.OutgoingMessage{
				Channel: chatMessage.Channel,
				User:    chatMessage.User,
				Body:    "Сабдей ещё не запущен SMOrc"},
			Command:    "voteCommand",
			RedirectTo: chatMessage.User})
		return
	}
	if subday.SubsOnly == true && chatMessage.IsSub == false {
		ircClient.SendDebounced(models.OutgoingDebouncedMessage{
			Message: models.OutgoingMessage{
				Channel: chatMessage.Channel,
				User:    chatMessage.User,
				Body:    "Ты не саб SMOrc"},
			Command:    "voteCommand",
			RedirectTo: chatMessage.User})
		return
	}
	repos.VoteForSubday(&chatMessage.User, &chatMessage.UserID, &subday.ID, &game)
}
