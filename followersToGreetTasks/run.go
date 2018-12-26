package followersToGreetSchedule

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/twitchIRCClient"
	"github.com/khades/servbot/userResolve"
	"sync"
	"time"
)

func Init(channelInfoService *channelInfo.Service,
	followersToGreetService *followersToGreet.Service,
	subAlertService *subAlert.Service,
	userResolveService *userResolve.Service,
	twitchIRCClient *twitchIRCClient.TwitchIRCClient,
	wg *sync.WaitGroup,
) *time.Ticker {
	ticker := time.NewTicker(time.Second * 20)

	service := Service{
		channelInfoService,
		followersToGreetService,
		subAlertService,
		userResolveService,
		twitchIRCClient,
	}
	go func(wg *sync.WaitGroup) {
		for {
			<-ticker.C
			wg.Add(1)
			service.AnnounceFollowers()
			wg.Done()
		}
	}(wg)

	return ticker
}