package autoMessageTasks

import (
	"html"
	"strings"

	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRC"

	"github.com/cbroglie/mustache"
)

type Service struct {
	channelInfoService *channelInfo.Service
	automessageService *autoMessage.Service
	twitchIRCClient    *twitchIRC.Client
}

func (service *Service) processMessage(message *autoMessage.AutoMessage) {

	channelInfo, error := service.channelInfoService.Get(&message.ChannelID)
	if error != nil {
		return
	}
	service.automessageService.ResetThreshold(message)

	compiledTemplate, templateError := mustache.ParseString(message.Message)
	if templateError != nil {
		return
	}
	compiledMessage, compiledMessageError := compiledTemplate.Render(channelInfo)
	if compiledMessageError != nil || strings.TrimSpace(compiledMessage) == "" {
		return
	}
	service.twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
		Channel: channelInfo.Channel,
		Body:    html.UnescapeString(compiledMessage)})

}

// Send checks if there any expired messages in chat and sends them
func (service *Service) Send() {
	messages, error := service.automessageService.ListCurrent()
	if error != nil {
		return
	}
	for _, message := range messages {
		service.processMessage(&message)
	}
}
