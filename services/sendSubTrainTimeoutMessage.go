package services

import (
	"html"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func SendSubTrainTimeoutMessage() {
	channels, error := repos.GetChannelsWithExpiredSubtrain()

	if error != nil {
		return
	}
	for _, channel := range channels {
		compiledMessage, compiledMessageError := mustache.Render(channel.SubTrain.TimeoutTemplate, channel.SubTrain)
		if compiledMessageError != nil || strings.TrimSpace(compiledMessage) != "" {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channel.Channel,
				Body:    html.UnescapeString(compiledMessage)})
		}
		eventbus.EventBus.Publish(eventbus.Subtrain(&channel.ChannelID), "expired")
		repos.ResetSubtrainCounter(&channel)
	}
}
