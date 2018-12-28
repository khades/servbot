package streamStatusTasks

import (
	"time"

	"github.com/khades/servbot/streamStatus"
	"github.com/sirupsen/logrus"
)

func Run(streamStatusService *streamStatus.Service) *time.Ticker {
	streamStatusService.UpdateFromTwitch()

	ticker := time.NewTicker(time.Second * 60)

	go func() {
		logger := logrus.WithFields(logrus.Fields{
			"package": "streamStatusTasks",
			"action":  "CheckStreamStatuses"})
		for range ticker.C {

			logger.Debug("Starting streams check")
			
			error := streamStatusService.UpdateFromTwitch()
			if error != nil {
				logger.Debugf("Error while updating streans: %s", error.Error())
			}
		}
	}()

	return ticker
}
