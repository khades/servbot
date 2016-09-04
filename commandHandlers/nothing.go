package commandHandlers

import (
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
)

// Nothing does Nothing
func Nothing(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	return
}
