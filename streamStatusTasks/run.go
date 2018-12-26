package streamStatusSchedule

import (
	"github.com/khades/servbot/streamStatus"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func Init(streamStatusService *streamStatus.Service,
	wg *sync.WaitGroup)  *time.Ticker {
	streamStatusService.UpdateFromTwitch()

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

			error := streamStatusService.UpdateFromTwitch()
			if error != nil {
				logger.Debugf("Error while updating streans: %s", error.Error())
			}
			wg.Done()
		}
	}(wg)

	return statusCheckerTicker
}
