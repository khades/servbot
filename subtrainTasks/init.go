package subtrainTasks

import (
	"time"

	"github.com/asaskevich/EventBus"

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
		for range ticker.C {
			service.Announce()
		}
	}()
	return ticker
}
