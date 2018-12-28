package autoMessageTasks

import (
	"time"

	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchIRC"
)

func Run(channelInfoService *channelInfo.Service,
	automessageService *autoMessage.Service,
	twitchIRCClient *twitchIRC.Client) *time.Ticker {
	ticker := time.NewTicker(time.Second * 20)

	service := Service{
		channelInfoService,
		automessageService,
		twitchIRCClient,
	}
	go func() {
		for range ticker.C{
			service.Send()
		}
	}()
	return ticker
}
