package bot

import (
	"strings"

	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

func modHandler(channel *string, mods []string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "bot",
		"feature": "mod",
		"action":  "modHandler"})
	filteredMods := []string{}
	for _, mod := range mods {
		if mod != "" {
			filteredMods = append(filteredMods, mod)
		}

	}
	logger.Debugf("Mods on channel %s: %s", *channel, strings.Join(filteredMods, ", "))
	users, error := repos.GetUsersID(filteredMods)
	if error != nil {
		logger.Debugf("GetUserID error: %s", error.Error())
		return
	}
	values, error := repos.GetUsersID([]string{*channel})
	channelID := (*values)[*channel]

	if error != nil {
		logger.Debugf("GetUserID error: %s", error.Error())
		return
	}
	if channelID == "" {
		logger.Debugf("ChannelID is not found for channel %s", *channel)

		return
	}
	userIDs := []string{}
	for _, id := range *users {
		userIDs = append(userIDs, id)
	}
	logger.Debugf("Mods IDs on channel %s: %s", *channel, strings.Join(userIDs, ", "))

	repos.PushMods(&channelID, userIDs)
}
