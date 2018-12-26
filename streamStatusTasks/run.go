package streamStatusTasks

import (
	"github.com/khades/servbot/streamStatus"
	"github.com/sirupsen/logrus"
	"time"
)

func Run(streamStatusService *streamStatus.Service)  *time.Ticker {
	streamStatusService.UpdateFromTwitch()

	statusCheckerTicker := time.NewTicker(time.Second * 60)

	go func() {
		for {
			<-statusCheckerTicker.C
			logger := logrus.WithFields(logrus.Fields{
				"package": "services",
				"feature": "streamstatus",
				"action":  "CheckStreamStatuses"})
			logger.Debug("Starting streams check")

			error := streamStatusService.UpdateFromTwitch()
			if error != nil {
				logger.Debugf("Error while updating streans: %s", error.Error())
			}
		}
	}()

	return statusCheckerTicker
}
