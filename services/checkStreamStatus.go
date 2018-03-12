package services

import (
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

// CheckStreamStatuses gets all stream statuses from all channels, and processes them
func CheckStreamStatuses() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "streamstatus",
		"action":  "CheckStreamStatuses"})
	logger.Debug("Starting streams check")

	error := repos.UpdateStreamStatuses()
	if error != nil {
		logger.Debugf("Error while updating streans: %s", error.Error())
	}
}
