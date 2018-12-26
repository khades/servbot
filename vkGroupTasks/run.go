package vkGroupTasks

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchIRC"
	"time"
)

func Run(config *config.Config, channelInfoService *channelInfo.Service, twitchIRCClient *twitchIRC.Client) *time.Ticker {
	ticker := time.NewTicker(time.Second * 60)
	service := Service{config, channelInfoService, twitchIRCClient}
	go func() {
		for {
			<-ticker.C
			service.Check()
		}
	}()
	return ticker
}
