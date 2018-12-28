package vkGroupTasks

import (
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchIRC"
	"github.com/sirupsen/logrus"
)

func Run(config *config.Config, channelInfoService *channelInfo.Service, twitchIRCClient *twitchIRC.Client) *time.Ticker {
	logger := logrus.WithFields(logrus.Fields{
		"package": "vkGroupTasks",
		"action":  "Check"})

	logger.Debugf("Checking VK")

	if config.VkClientKey == "" {
		logger.Infof("VK key is not set")
		return nil
	}

	ticker := time.NewTicker(time.Second * 60)
	service := Service{config, channelInfoService, twitchIRCClient}
	go func() {
		for range ticker.C {
			service.Check()
		}
	}()
	return ticker
}
