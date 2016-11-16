package bot

import (
	"log"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func sendSubMessage(channel *string, user *string) {
	subAlert, error := repos.GetSubAlert(channel)
	log.Println(*subAlert)
	if error == nil && subAlert.Enabled == true && subAlert.SubMessage != "" {
		IrcClientInstance.SendPublic(&models.OutgoingMessage{
			Body:    subAlert.SubMessage,
			Channel: *channel,
			User:    *user})
	}
}
