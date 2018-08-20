package services

import (
	"strings"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

// AnnounceFollowers announces all new followers on channels
func AnnounceFollowers() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "followers",
		"action":  "AnnounceFollowers"})
	channelFollowers, error := repos.GetFollowersToGreet()
	if error != nil {
		logger.Debug("Nothing to process")
		return
	}
	for _, channel := range channelFollowers {
		processOneChannel(channel)
	}

}

func processOneChannel(channel models.FollowersToGreet) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "followers",
		"action":  "processOneChannel"})
	defer repos.ResetFollowersToGreetOnChannel(&channel.ChannelID)
	logger.Debugf("Processing channel %s", channel.ChannelID)
	alertInfo, alertError := repos.GetSubAlert(&channel.ChannelID)
	if alertError != nil {
		logger.Debugf("No alert for channel %s", channel.ChannelID)

		return
	}
	channelInfo, channelInfoError := repos.GetChannelInfo(&channel.ChannelID)
	if channelInfoError != nil {
		logger.Debugf("No channelInfo for channel %s", channel.ChannelID)

		return
	}
	followers := []string{}
	followersMap, followersError := repos.GetUsernames(channel.Followers)
	if followersError != nil {
		logger.Debugf("Followers resolve failed for channel %s", channel.ChannelID)

		return
	}
	for _, follower := range *followersMap {
		if (strings.TrimSpace(follower) != "") {
			followers = append(followers, follower)
		}
	}
	
	followersString := strings.TrimSpace(strings.Join(followers, " @"))
	if channelInfoError == nil && channelInfo.Enabled == true && alertInfo.Enabled == true && followersString != "" && alertInfo.FollowerMessage != "" {
		bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
			Channel: channelInfo.Channel,
			Body:    "@" + followersString + " " + alertInfo.FollowerMessage})
	}
}
