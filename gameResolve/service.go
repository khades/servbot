package gameResolve

import (
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchAPI"
	"github.com/sirupsen/logrus"
)

type Service struct {
	collection         *mgo.Collection
	twitchAPIClient    *twitchAPI.Client
	channelInfoService *channelInfo.Service
	gamesToProcess     []string
}

// Get returns game by its game id,
func (service *Service) Get(gameID *string) (string, bool) {
	if strings.TrimSpace(*gameID) == "" {
		return "", false
	}
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "game",
		"action":  "GetGamesByID"})
	var result Game
	error := service.collection.Find(
		bson.M{"gameid": *gameID,
			"date": bson.M{"$gte": time.Now().Add(-14 * 24 * time.Hour)}}).One(&result)
	if error != nil || result.Game == "" {
		service.gamesToProcess = append(service.gamesToProcess, *gameID)
		logger.Debugf("Can't find game %s in Database: %s", *gameID, error.Error())
		return "", false
	}
	return result.Game, true
}

func (service *Service) set(gameID *string, game *string) {
	service.collection.Upsert(
		bson.M{"gameid": *gameID},
		bson.M{"$set": bson.M{"game": *game, "date": time.Now()}})
}

// Fetch processes games that were not found in database previously
func (service *Service) Fetch() error {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "twitchGames",
		"action":  "Fetch"})

	logger.Debugf("%d games were not found in database", len(service.gamesToProcess))
	if len(service.gamesToProcess) == 0 {
		return nil
	}
	games, error := service.twitchAPIClient.GetGamesByID(service.gamesToProcess)
	if error != nil {
		return error
	}
	for _, game := range games {
		service.set(&game.GameID, &game.Game)
		service.channelInfoService.UpdateGamesByID(&game.GameID, &game.Game)
	}
	service.gamesToProcess = []string{}
	return nil
}
