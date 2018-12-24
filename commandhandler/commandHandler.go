package commandhandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
)

// CommandHandle is interface for functions that will handle stuff
type CommandHandle func(channelInfo *channelInfo.ChannelInfo, chatMessage *chatMessage.ChatMessage, chatCommand models.ChatCommand)
