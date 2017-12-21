package services

import (
	"html"
	"strings"


	"github.com/hoisie/mustache"
	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func SendSubTrainTimeoutMessage() {
	channels, error := repos.GetChannelsWithExpiredSubtrain()

	if error != nil {
		return
	}
	for _, channel := range *channels {
		compiledMessage := mustache.Render(channel.SubTrain.TimeoutTemplate, channel.SubTrain)
		if strings.TrimSpace(compiledMessage) != "" {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channel.Channel,
				Body:    html.UnescapeString(compiledMessage)})
		}
		repos.ResetSubtrainCounter(&channel)
	}
}
