package twitchIRCHandler

import (
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRCClient"
)

func (service *TwitchIRCHandler) sendSubMessage(client *twitchIRCClient.TwitchIRCClient, channelInfo *channelInfo.ChannelInfo, user *string, subPlan *string) {
	subAlert, error := service.subAlertService.Get(&channelInfo.ChannelID)

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

		template = template + " " + strings.TrimSpace(subtrainAdditionalString)
		service.channelInfoService.IncrementSubtrainCounterByChannelID(&channelInfo.ChannelID, user)
	}

	if template != "" {
		client.SendPublic(&twitchIRCClient.OutgoingMessage{
			Body:    template,
			Channel: channelInfo.Channel,
			User:    *user})
	}
}
