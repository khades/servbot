package bot

import (
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func sendSubMessage(channel *string, channelID *string, user *string, subPlan *string) {
	subAlert, error := repos.GetSubAlert(channelID)

	if error != nil || subAlert.Enabled == false {
		return
	}
	template := subAlert.SubFiveMessage
	switch *subPlan {
	case "Prime":
		{
			if subAlert.SubPrimeMessage != "" {
				template = subAlert.SubPrimeMessage
			}
		}
	case "2000":
		{
			if subAlert.SubTenMessage != "" {
				template = subAlert.SubTenMessage
			}
		}
	case "3000":
		{
			if subAlert.SubTwentyFiveMessage != "" {
				template = subAlert.SubTwentyFiveMessage
			}
		}
	}
	if template != "" {
		IrcClientInstance.SendPublic(&models.OutgoingMessage{
			Body:    template,
			Channel: *channel,
			User:    *user})
	}
}
