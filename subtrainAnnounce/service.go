package subtrainAnnounce

import (
	"github.com/asaskevich/EventBus"
	"html"
	"strings"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRCClient"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/eventbus"
)

type Service struct {
	channelInfoService *channelInfo.Service
	twitchIRCClient    *twitchIRCClient.TwitchIRCClient
	eventBus 		   EventBus.Bus
}

// Announce gets all expired subtrains and sends expiration message in channels
func (service *Service) Announce() {
	channels, error := service.channelInfoService.GetChannelsWithExpiredSubtrain()

	if error != nil {
		return
	}
	for _, channel := range channels {
		compiledMessage, compiledMessageError := mustache.Render(channel.SubTrain.TimeoutTemplate, channel.SubTrain)
		if compiledMessageError != nil || strings.TrimSpace(compiledMessage) != "" {
			service.twitchIRCClient.SendPublic(&twitchIRCClient.OutgoingMessage{
				Channel: channel.Channel,
				Body:    html.UnescapeString(compiledMessage)})
		}
		service.eventBus.Publish(eventbus.Subtrain(&channel.ChannelID), "expired")
		service.channelInfoService.ResetSubtrainCounter(&channel)
	}
	channels, error = service.channelInfoService.GetChannelsWithSubtrainNotification()
	if error != nil {
		return
	}
	for _, channel := range channels {
		compiledMessage, compiledMessageError := mustache.Render(channel.SubTrain.NotificationTemplate, channel.SubTrain)
		if compiledMessageError != nil && strings.TrimSpace(compiledMessage) != "" {
			service.twitchIRCClient.SendPublic(&twitchIRCClient.OutgoingMessage{
				Channel: channel.Channel,
				Body:    html.UnescapeString(compiledMessage)})
		}
		service.channelInfoService.SetSubtrainNotificationShown(&channel)
	}
}
