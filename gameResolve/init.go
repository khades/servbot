package gameResolve

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchAPIClient"
)

var gamesCollection = "games"

func Init(db *mgo.Database,
	twitchAPIClient *twitchAPIClient.TwitchAPIClient,
	channelInfoService *channelInfo.Service,
	group *sync.WaitGroup) (*Service, *time.Ticker) {

	collection := db.C(gamesCollection)
	collection.EnsureIndex(mgo.Index{
		Key: []string{"gameid"}})

	service := &Service{
		collection:         collection,
		twitchAPIClient:          twitchAPIClient,
		channelInfoService: channelInfoService,
		gamesToProcess:     []string{},
	}

	gamesCheckerTicker := time.NewTicker(time.Second * 30)

	go func(wg *sync.WaitGroup) {
		for {
			<-gamesCheckerTicker.C
			wg.Add(1)
			logger := logrus.WithFields(logrus.Fields{
				"package": "services",
				"feature": "twitchGames",
				"action":  "GetTwitchGames"})
			logger.Debug("Starting twitch games fetching")
			error := service.Fetch()
			if error != nil {
				logger.Debug("Twitch games fetching error: %s", error.Error())
			}
			wg.Done()
		}
	}(group)

	return service, gamesCheckerTicker
}
