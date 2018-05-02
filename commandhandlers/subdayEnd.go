package commandhandlers

// import (
// 	//"strings"
// 	//"time"

// 	"github.com/khades/servbot/ircClient"
// 	"github.com/khades/servbot/l10n"
// 	"github.com/khades/servbot/models"
// 	"github.com/khades/servbot/repos"
// )

// func subdayEnd(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {

// 	if chatMessage.IsMod == false {
// 		// ircClient.SendPublic(&models.OutgoingMessage{
// 		// 	Channel: chatMessage.Channel,
// 		// 	Body:    l10n.GetL10n(channelInfo.GetChannelLang()).SubdayEndYoureNotModerator,
// 		// 	User:    chatMessage.User})
// 		return
// 	}
// 	if channelInfo.SubdayIsActive == true {
// 		repos.CloseActiveSubday(&chatMessage.ChannelID)
// 		ircClient.SendPublic(&models.OutgoingMessage{
// 			Channel: chatMessage.Channel,
// 			Body: 	 l10n.GetL10n(channelInfo.GetChannelLang()).SubdayEndSuccess,
// 			User:    chatMessage.User})
// 		return
// 	}
// 	ircClient.SendPublic(&models.OutgoingMessage{
// 		Channel: chatMessage.Channel,
// 		Body:  	 l10n.GetL10n(channelInfo.GetChannelLang()).SubdayEndNothingToClose,
// 		User:    chatMessage.User})
// 	return

// }
