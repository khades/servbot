package configTasks

import (
	"time"

	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchAPI"
	"github.com/sirupsen/logrus"
)

type Service struct {
	twitchAPIClient *twitchAPI.Client
	config          *config.Config
}

func (service *Service) UpdateAPIKey(period time.Duration) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "configTasks",
		"action":  "UpdateAPIKey"})
	needsUpdate := service.config.NeedsAPIKey(period)
	if needsUpdate == false {
		logger.Debug("APIKey doesn't need to be updated")
		return
	}
	logger.Debug("APIKey does need to be updated, Processing.")

	newKeyStruct, error := service.twitchAPIClient.GetAPIKey()
	if error != nil {
		logger.Debugf("APIKey retrival error : %s", error.Error())
		// needs logging
		return
	}
	service.config.SaveApiKey(newKeyStruct.AccessToken, newKeyStruct.ExpiresIn)
}
