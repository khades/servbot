package repos

import (
	"time"

	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
	"github.com/globalsign/mgo/bson"
)

var gamesCollection = "games"
var gamesToProcess = []string{}

// GetGameByID returns game by its game id,
func getGameByID(gameID *string) (string, bool) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "game",
		"action":  "GetGamesByID"})
	var result models.Game
	error := db.C(gamesCollection).Find(
		bson.M{"gameid": *gameID,
			"date": bson.M{"$gte": time.Now().Add(-14 * 24 * time.Hour)}}).One(&result)
	if error != nil || result.Game == "" {
		gamesToProcess = append(gamesToProcess, *gameID)
		logger.Debugf("Can't find game in Database: %s", error.Error())
		return "", false
	}
	return result.Game, true
}

func setGameByID(gameID *string, game *string) {
	db.C(gamesCollection).Upsert(
		bson.M{"gameid": *gameID},
		bson.M{"$set": bson.M{"game": *game, "date": time.Now()}})
}


// FetchGamesFromTwitch processes games that were not found in database previously
func FetchGamesFromTwitch() error {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "twitchGames",
		"action":  "FetchGamesFromTwitch"})

	logger.Debugf("%d games were not found in database", len(gamesToProcess))
	if len(gamesToProcess) == 0 {
		return nil
	}
	games, error := getGamesByID(gamesToProcess)
	if error != nil {
		return error
	}
	for _, game := range games {
		setGameByID(&game.GameID, &game.Game)
		updateGamesByID(&game.GameID, &game.Game)
	}
	gamesToProcess = []string{}
	return nil
}