package followersToGreetTasks

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/twitchIRC"
	"github.com/khades/servbot/userResolve"
	"strings"

	"github.com/sirupsen/logrus"
)

type Service struct {
	channelInfoService *channelInfo.Service
	followersToGreetService *followersToGreet.Service
	subAlertService *subAlert.Service
	userResolveService *userResolve.Service
	twitchIRCClient *twitchIRC.Client
}
// AnnounceFollowers announces all new followers on channels
func (service *Service) AnnounceFollowers() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "followers",
		"action":  "AnnounceFollowers"})
	channelFollowers, error := service.followersToGreetService.List()
	if error != nil {
		logger.Debug("Nothing to process")
		return
	}
	for _, channel := range channelFollowers {
		service.processOneChannel(channel)
	}

}

func (service *Service) processOneChannel(channel followersToGreet.FollowersToGreet) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "followers",
		"action":  "processOneChannel"})
	defer service.followersToGreetService.Reset(&channel.ChannelID)
	logger.Debugf("Processing channel %s", channel.ChannelID)
	alertInfo, alertError := service.subAlertService.Get(&channel.ChannelID)
	if alertError != nil {
		logger.Debugf("No alert for channel %s", channel.ChannelID)

		return
	}
	channelInfo, channelInfoError := service.channelInfoService.Get(&channel.ChannelID)
	if channelInfoError != nil {
		logger.Debugf("No channelInfo for channel %s", channel.ChannelID)

		return
	}
	followers := []string{}
	logger.Debugf("Looking for followers names: ", strings.Join(channel.Followers, ", "))
	followersMap, followersError := service.userResolveService.GetUsernames(channel.Followers)
	if followersError != nil {
		logger.Debugf("Followers resolve failed for channel %s", channel.ChannelID)

		return
	}
	logger.Debugf("Search result: %+v", followersMap)

	for _, follower := range *followersMap {
		if (strings.TrimSpace(follower) != "") {
			followers = append(followers, follower)
		}
	}
	
	followersString := strings.TrimSpace(strings.Join(followers, " @"))
	if channelInfoError == nil && channelInfo.Enabled == true && alertInfo.Enabled == true && followersString != "@" && alertInfo.FollowerMessage != "" {
		service.twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
			Channel: channelInfo.Channel,
			Body:    "@" + followersString + " " + alertInfo.FollowerMessage})
	}
}
