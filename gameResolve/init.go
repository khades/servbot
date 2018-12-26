package gameResolve

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchAPI"
)

const collectionName  = "games"

func Init(db *mgo.Database,
	twitchAPIClient *twitchAPI.Client,
	channelInfoService *channelInfo.Service) (*Service, *time.Ticker) {

	collection := db.C(collectionName)
	collection.EnsureIndex(mgo.Index{
		Key: []string{"gameid"}})

	service := &Service{
		collection:         collection,
		twitchAPIClient:          twitchAPIClient,
		channelInfoService: channelInfoService,
		gamesToProcess:     []string{},
	}

	gamesCheckerTicker := time.NewTicker(time.Second * 30)

	go func() {
		for {
			<-gamesCheckerTicker.C

			logger := logrus.WithFields(logrus.Fields{
				"package": "services",
				"feature": "twitchGames",
				"action":  "GetTwitchGames"})
			logger.Debug("Starting twitch games fetching")
			error := service.Fetch()
			if error != nil {
				logger.Debug("Twitch games fetching error: %s", error.Error())
			}
		}
	}()

	return service, gamesCheckerTicker
}
