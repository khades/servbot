package streamStatus

import (
	"sync"
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/gameResolve"
	"github.com/khades/servbot/twitchAPIClient"
	"github.com/sirupsen/logrus"
)

func Init(config *config.Config,
	channelInfoService *channelInfo.Service,
	gameResolveService *gameResolve.Service,
	twitchAPIClient *twitchAPIClient.TwitchAPIClient,
	wg *sync.WaitGroup) (*Service, *time.Ticker) {
	service := Service{config: config, channelInfoService: channelInfoService, gameResolveService: gameResolveService, twitchAPIClient: twitchAPIClient}
	service.UpdateFromTwitch()

	statusCheckerTicker := time.NewTicker(time.Second * 60)

	go func(wg *sync.WaitGroup) {
		for {
			<-statusCheckerTicker.C
			wg.Add(1)
			logger := logrus.WithFields(logrus.Fields{
				"package": "services",
				"feature": "streamstatus",
				"action":  "CheckStreamStatuses"})
			logger.Debug("Starting streams check")

			error := service.UpdateFromTwitch()
			if error != nil {
				logger.Debugf("Error while updating streans: %s", error.Error())
			}
			wg.Done()
		}
	}(wg)

	return &service, statusCheckerTicker
}
