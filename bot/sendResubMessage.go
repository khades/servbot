package bot

import (
	"log"
	"strings"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func sendResubMessage(channel *string, user *string, resubCount *int) {
	subAlert, error := repos.GetSubAlert(channel)
	log.Println(*subAlert)
	if error == nil && subAlert.Enabled == true && subAlert.ResubMessage != "" {
		template, error := repos.ResubTemplateCache.Get(subAlert)
		if error == nil {
			resubInfo := models.ResubInfo{RepeatedBody: strings.Repeat(subAlert.RepeatBody+" ", *resubCount), ResubCount: *resubCount}
			compiledMessage := template.Render(resubInfo)
			if compiledMessage != "" {
				IrcClientInstance.SendPublic(&models.OutgoingMessage{
					Body:    compiledMessage,
					Channel: *channel,
					User:    *user})
			}
		}

	}
}
