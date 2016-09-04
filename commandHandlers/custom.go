package commandHandlers

import (
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
)

// Custom handler does template job
func Custom(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	found, template := Template.get(chatMessage.Channel, chatCommand.Command)
	if found {
		message := template.Render(chatMessage)

		ircClient.SendDebounced(models.OutgoingDebouncedMessage{
			Message: models.OutgoingMessage{
				Channel: chatMessage.Channel,
				User:    chatMessage.User,
				Body:    message},
			Command:    chatCommand.Command,
			Redirected: false})
	}
}
