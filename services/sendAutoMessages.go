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
		Channel: channelInfo.Channel,
		Body:    html.UnescapeString(compiledMessage)})

}

// SendAutoMessages checks if there any expired messages in chat and sends them
func SendAutoMessages() {
	messages, error := repos.GetCurrentAutoMessages()
	if error != nil {
		return
	}
	for _, message := range messages {
		processMessage(&message)
	}
}
