package twitchIRCHandler

import (
	"strings"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRC"

	"github.com/cbroglie/mustache"
)

func (service *TwitchIRCHandler) sendResubMessage(client *twitchIRC.Client, channelInfo *channelInfo.ChannelInfo, user *string, resubCount *int, subPlan *string) {
	subAlert, error := service.subAlertService.Get(&channelInfo.ChannelID)

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
		resubInfo := ResubInfo{Smiles: strings.Repeat(smile+" ", *resubCount), ResubCount: *resubCount}
		compiledMessage, _ := compiledTemplate.Render(resubInfo)
		if channelInfo.SubTrain.Enabled && channelInfo.SubTrain.OnlyNewSubs == false {
			localSubtrain := channelInfo.SubTrain
			localSubtrain.CurrentStreak = localSubtrain.CurrentStreak + 1
			subtrainAdditionalString, _ := mustache.Render(channelInfo.SubTrain.AppendTemplate, localSubtrain)
			compiledMessage = compiledMessage + " " + strings.TrimSpace(subtrainAdditionalString)
			service.channelInfoService.IncrementSubtrainCounterByChannelID(&channelInfo.ChannelID, user)
		}

		if compiledMessage != "" {
			client.SendPublic(&twitchIRC.OutgoingMessage{
				Body:    compiledMessage,
				Channel: channelInfo.Channel,
				User:    *user})
		}
	}

}
