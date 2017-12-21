package bot

import (
	"strings"

	"github.com/hoisie/mustache"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func sendResubMessage(channel *string, channelID *string, user *string, resubCount *int, subPlan *string) {
	subAlert, error := repos.GetSubAlert(channelID)

	if error != nil || subAlert.Enabled == false {
		return
	}
	template := subAlert.ResubFiveMessage
	smile := subAlert.ResubFiveSmile
	switch *subPlan {
	case "Prime":
		{
			if subAlert.ResubPrimeMessage != "" {
				template = subAlert.ResubPrimeMessage
				smile = subAlert.ResubPrimeSmile
			}
		}
	case "2000":
		{
			if subAlert.ResubTenMessage != "" {
				template = subAlert.ResubTenMessage
				smile = subAlert.ResubTenSmile
			}
		}
	case "3000":
		{
			if subAlert.ResubTwentyFiveMessage != "" {
				template = subAlert.ResubTwentyFiveMessage
				smile = subAlert.ResubTwentyFiveSmile
			}
		}
	}
	compiledTemplate, error := mustache.ParseString(template)
	if error == nil {
		resubInfo := models.ResubInfo{Smiles: strings.Repeat(smile+" ", *resubCount), ResubCount: *resubCount}
		compiledMessage := compiledTemplate.Render(resubInfo)
		channelInfo, channelInfoError := repos.GetChannelInfo(channelID)
		if channelInfoError == nil && channelInfo.SubTrain.Enabled {
			localSubtrain := channelInfo.SubTrain
			localSubtrain.CurrentStreak = localSubtrain.CurrentStreak + 1
			compiledMessage = compiledMessage + mustache.Render(channelInfo.SubTrain.AppendTemplate, localSubtrain)
		}
		if compiledMessage != "" {
			IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Body:    compiledMessage,
				Channel: *channel,
				User:    *user})
		}
	}

}
