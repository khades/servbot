package webhook

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchAPIClient"
	"sync"
	"time"
)

var webhooklibrary = "webhooklibrary"

func Init(db *mgo.Database, channelInfoService *channelInfo.Service, twitchAPIService *twitchAPIClient.TwitchAPIClient, wg *sync.WaitGroup) (*Service, *time.Ticker) {
	collection := db.C(webhooklibrary)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "topic"}})
		
	service := &Service{
		collection: collection,
		channelInfoService: channelInfoService,
		twitchAPIService: twitchAPIService,
	}

	webhookTimer := time.NewTicker(time.Minute * 15)
	service.Subscribe(time.Minute * 15)
	go func(wg *sync.WaitGroup) {
		for {
			<-webhookTimer.C
			wg.Add(1)
			service.Subscribe(time.Minute * 15)
			wg.Done()
		}
	}(wg)

	return service, webhookTimer
}
