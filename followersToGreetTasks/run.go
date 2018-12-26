package followersToGreetTasks

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/twitchIRC"
	"github.com/khades/servbot/userResolve"
	"time"
)

func Run(channelInfoService *channelInfo.Service,
	followersToGreetService *followersToGreet.Service,
	subAlertService *subAlert.Service,
	userResolveService *userResolve.Service,
	twitchIRCClient *twitchIRC.Client) *time.Ticker {
	ticker := time.NewTicker(time.Second * 20)

	service := Service{
		channelInfoService,
		followersToGreetService,
		subAlertService,
		userResolveService,
		twitchIRCClient,
	}
	go func() {
		for {
			<-ticker.C
			service.AnnounceFollowers()
		}
	}()

	return ticker
}