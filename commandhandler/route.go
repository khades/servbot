package commandhandler

import (
	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/followers"
	"github.com/khades/servbot/twitchIRCClient"
	"github.com/khades/servbot/songRequest"
	"github.com/khades/servbot/subday"
	"github.com/khades/servbot/template"
)

// RouterStruct is struct for handling command to handler
type CommandHandler struct {
	templateService *template.Service
	subdayService *subday.Service
	followersService *followers.Service
	songRequestService *songRequest.Service
	twitchIRCClient       *bot.TwitchIRCClient
}

// Go returns from router to work with
func (commandHandler CommandHandler) Go(command string) CommandHandle {
	if command == "new" {
		return commandHandler.newCommand
	}
	if command == "alias" {
		return commandHandler.aliasCommand
	}
	if command == "subdaynew" {
		return commandHandler.subdayNew
	}
	return commandHandler.custom
}


