package commandHandlers

import (
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
)

// CommandHandler is interface for functions that will handle stuff
type CommandHandler func(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient)
