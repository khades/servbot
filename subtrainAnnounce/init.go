package subtrainAnnounce

import (
	"github.com/asaskevich/EventBus"
	"sync"
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRCClient"
)

func Init(channelInfoService *channelInfo.Service,
	twitchIRCClient *twitchIRCClient.TwitchIRCClient,
	eventBus EventBus.Bus,
	wg *sync.WaitGroup) *time.Ticker {
	subtrainAnnounceTicker := time.NewTicker(time.Second * 10)

	service := Service{
		channelInfoService,
		twitchIRCClient,
		eventBus,
	}
	go func(wg *sync.WaitGroup) {
		for {
			<-subtrainAnnounceTicker.C
			wg.Add(1)
			service.Announce()
			wg.Done()
		}
	}(wg)
	return subtrainAnnounceTicker
}
