package commandhandlers

import (
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
)

// Nothing does nothing, it is stub
func Nothing(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage,  chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	return
}
