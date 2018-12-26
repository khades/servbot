package subtrainTasks

import (
	"github.com/asaskevich/EventBus"
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRC"
)

func Run(channelInfoService *channelInfo.Service,
	twitchIRCClient *twitchIRC.Client,
	eventBus EventBus.Bus) *time.Ticker {
	ticker := time.NewTicker(time.Second * 10)

	service := Service{
		channelInfoService,
		twitchIRCClient,
		eventBus,
	}
	go func() {
		for {
			<-ticker.C
			service.Announce()
		}
	}()
	return ticker
}
