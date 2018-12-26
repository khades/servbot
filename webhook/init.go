package webhook

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/twitchAPIClient"
)

var webhooklibrary = "webhooklibrary"

func Init(db *mgo.Database, channelInfoService *channelInfo.Service, twitchAPIService *twitchAPIClient.TwitchAPIClient) *Service {
	collection := db.C(webhooklibrary)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "topic"}})

	service := &Service{
		collection:         collection,
		channelInfoService: channelInfoService,
		twitchAPIService:   twitchAPIService,
	}

	return service
}
