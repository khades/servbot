package autoMessageAnnounce

import (
	"sync"
	"time"

	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRCClient"
)

func Init(channelInfoService *channelInfo.Service,
	automessageService *autoMessage.Service,
	twitchIRCClient *twitchIRCClient.TwitchIRCClient,
	wg *sync.WaitGroup) *time.Ticker {
	automessageTicker := time.NewTicker(time.Second * 20)

	service := Service{
		channelInfoService,
		automessageService,
		twitchIRCClient,
	}
	go func(wg *sync.WaitGroup) {
		for {
			<-automessageTicker.C
			wg.Add(1)
			service.Send()
			wg.Done()
		}
	}(wg)
	return automessageTicker;
}
