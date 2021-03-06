package services

import (
	"html"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func processMessage(message *models.AutoMessage) {
	channel, error := repos.GetUsernameByID(&message.ChannelID)
	if error != nil {
		return
	}
	channelInfo, error := repos.GetChannelInfo(&message.ChannelID)
	if error != nil {
		return
	}
	repos.ResetAutoMessageThreshold(message)

	compiledTemplate, templateError := mustache.ParseString(message.Message)
	if templateError != nil {
		return
	}
	compiledMessage, compiledMessageError := compiledTemplate.Render(channelInfo)
	if compiledMessageError != nil || strings.TrimSpace(compiledMessage) == "" {
		return
	}
	bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
		Channel: *channel,
		Body:    html.UnescapeString(compiledMessage)})

}

func SendAutoMessages() {
	messages, error := repos.GetCurrentAutoMessages()
	if error != nil {
		return
	}
	for _, message := range *messages {
		processMessage(&message)
	}
}
