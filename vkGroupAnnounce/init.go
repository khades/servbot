package vkGroupAnnounce

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchIRCClient"
	"sync"
	"time"
)

func Init(config *config.Config, channelInfoService *channelInfo.Service, twitchIRCClient *twitchIRCClient.TwitchIRCClient, wg *sync.WaitGroup) *time.Ticker {
	vkTimer := time.NewTicker(time.Second * 60)
	service := Service{config, channelInfoService, twitchIRCClient}
	go func(wg *sync.WaitGroup) {
		for {

			<-vkTimer.C
			wg.Add(1)
			service.Check()
			wg.Done()
		}
	}(wg)
	return vkTimer
}
