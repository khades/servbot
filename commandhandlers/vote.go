package commandhandlers

import (
	"unicode/utf8"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// Vote handlers processes incoming subday vote variants
func Vote(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	if utf8.RuneCountInString(chatMessage.MessageBody) < 2 {
		return
	}

	game := chatMessage.MessageBody[1:]
	subday, subdayError := repos.GetLastActiveSubday(&chatMessage.ChannelID)
	if subdayError != nil {

		ircClient.SendDebounced(models.OutgoingDebouncedMessage{Message: models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayVoteNoActiveSubday,
			User:    chatMessage.User},
			RedirectTo: chatMessage.User,
			Command:    "%vote"})
		return
	}
	if subday.SubsOnly == true && chatMessage.IsSub == false {
		ircClient.SendDebounced(models.OutgoingDebouncedMessage{Message: models.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayVoteYouReNotSub,
			User:    chatMessage.User},
			RedirectTo: chatMessage.User,
			Command:    "%vote"})
		return
	}
	repos.VoteForSubday(&chatMessage.User, &chatMessage.UserID, chatMessage.IsSub, &subday.ID, &game)
}
