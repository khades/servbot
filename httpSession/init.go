package httpSession

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/twitchAPIClient"
)

var httpsessionCollection = "httpSessions"

func Init(db *mgo.Database, twitchAPIClient *twitchAPIClient.TwitchAPIClient) *Service {
	collection := db.C(httpsessionCollection)

	collection.EnsureIndex(mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: 15 * time.Minute})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"key"}})

	return &Service{collection: collection, twitchAPIClient: twitchAPIClient}
}
