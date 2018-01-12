package services

import (
	"html"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func SendSubTrainNotification() {
	channels, error := repos.GetChannelsWithSubtrainNotification()
	if error != nil {
		return
	}
	for _, channel := range *channels {
		compiledMessage, compiledMessageError := mustache.Render(channel.SubTrain.NotificationTemplate, channel.SubTrain)
		if compiledMessageError != nil && strings.TrimSpace(compiledMessage) != "" {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channel.Channel,
				Body:    html.UnescapeString(compiledMessage)})
		}
		repos.SetSubtrainNotificationShown(&channel)
	}
}
