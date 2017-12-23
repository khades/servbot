package bot

import (
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"strings"
	"github.com/hoisie/mustache"
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
	channelInfo, channelInfoError := repos.GetChannelInfo(channelID)
	if channelInfoError == nil && channelInfo.SubTrain.Enabled {
		localSubtrain := channelInfo.SubTrain
		localSubtrain.CurrentStreak = localSubtrain.CurrentStreak + 1
		template = template +" " + strings.TrimSpace(mustache.Render(channelInfo.SubTrain.AppendTemplate, localSubtrain))
		repos.IncrementSubtrainCounterByChannelID(channelID, user)
	}
	
	if template != "" {
		IrcClientInstance.SendPublic(&models.OutgoingMessage{
			Body:    template,
			Channel: *channel,
			User:    *user})
	}
}
