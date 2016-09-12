package commandHandlers

import (
	"log"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// Custom handler checks if input command has template and then fills it in with mustache templating and sends to a specified/user
func Custom(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	template := Template.get(chatMessage.Channel, chatCommand.Command)
	if template != nil {
		redirected := false
		values, _ := repos.GetChannelInfo(chatMessage.Channel)
		message, templateError := template.Render(values)

		log.Println(templateError)
		log.Println(message)
		if templateError != nil {
			message = "Ошибка в шаблоне команды, обратитесь к модератору etmSad"
		}
		user := chatMessage.User
		if chatCommand.Body != "" {
			user = chatCommand.Body
			redirected = true
		}
		ircClient.SendDebounced(models.OutgoingDebouncedMessage{
			Message: models.OutgoingMessage{
				Channel: chatMessage.Channel,
				User:    user,
				Body:    message},
			Command:    chatCommand.Command,
			Redirected: redirected})
	}
}
