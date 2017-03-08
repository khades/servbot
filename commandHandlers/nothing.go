package commandHandlers

import (
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
)

// Nothing does nothing, it is stub
func Nothing(online bool, chatMessage *models.ChatMessage,  chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	return
}
