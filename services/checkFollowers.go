package services

import (
	"strings"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// CheckChannelsFollowers process followers of all channels on that instance of bot, that code will be deprecated after webhooks will be done
func CheckChannelsFollowers() {
	channelFollowers, error := repos.CheckFollowers()
	if error != nil {
		return
	}
	for _, channel := range channelFollowers{
		alertInfo, alertError := repos.GetSubAlert(&channel.ChannelID)

		if alertError == nil && alertInfo.Enabled == true && alertInfo.FollowerMessage != "" {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channel.Channel,
				Body:    "@" + strings.Join(channel.Followers, " @") + " " + alertInfo.FollowerMessage})
		}
	}

}
