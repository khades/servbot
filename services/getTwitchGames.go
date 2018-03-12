package services

import (
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

func GetTwitchGames() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "twitchGames",
		"action":  "GetTwitchGames"})
	logger.Debug("Starting twitch games fetching")
	error := repos.FetchGamesFromTwitch()
	if error != nil {
		logger.Debug("Twitch games fetching error: %s", error.Error())
	}
}
