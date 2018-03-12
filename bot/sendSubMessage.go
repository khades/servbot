package bot

import (
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"strings"
	"github.com/cbroglie/mustache"
)

func sendSubMessage(channelInfo *models.ChannelInfo,user *string, subPlan *string) {
	subAlert, error := repos.GetSubAlert(&channelInfo.ChannelID)

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
	if channelInfo.SubTrain.Enabled {
		localSubtrain := channelInfo.SubTrain
		localSubtrain.CurrentStreak = localSubtrain.CurrentStreak + 1
		subtrainAdditionalString, _ := mustache.Render(channelInfo.SubTrain.AppendTemplate, localSubtrain)

		template = template +" " + strings.TrimSpace(subtrainAdditionalString)
		repos.IncrementSubtrainCounterByChannelID(&channelInfo.ChannelID, user)
	}
	
	if template != "" {
		IrcClientInstance.SendPublic(&models.OutgoingMessage{
			Body:    template,
			Channel: channelInfo.Channel,
			User:    *user})
	}
}
