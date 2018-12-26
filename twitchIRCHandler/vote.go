package twitchIRCHandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/twitchIRC"
	"unicode/utf8"

	"github.com/khades/servbot/l10n"
)

// Vote handlers processes incoming subday vote variants
func (service *TwitchIRCHandler) vote(client *twitchIRC.Client, channelInfo *channelInfo.ChannelInfo, chatMessage *chatMessage.ChatMessage) {
	if utf8.RuneCountInString(chatMessage.MessageBody) < 2 {
		return
	}

	game := chatMessage.MessageBody[1:]
	subday, subdayError := service.subdayService.GetActive(&chatMessage.ChannelID)
	if subdayError != nil {

		client.SendDebounced(twitchIRC.OutgoingDebouncedMessage{Message: twitchIRC.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayVoteNoActiveSubday,
			User:    chatMessage.User},
			RedirectTo: chatMessage.User,
			Command:    "%vote"})
		return
	}
	if subday.SubsOnly == true && chatMessage.IsSub == false {
		client.SendDebounced(twitchIRC.OutgoingDebouncedMessage{Message: twitchIRC.OutgoingMessage{
			Channel: chatMessage.Channel,
			Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayVoteYouReNotSub,
			User:    chatMessage.User},
			RedirectTo: chatMessage.User,
			Command:    "%vote"})
		return
	}
	service.subdayService.Vote(&chatMessage.User, &chatMessage.UserID, chatMessage.IsSub, &subday.ID, &game)
}
