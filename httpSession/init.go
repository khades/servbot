package httpSession

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/twitchAPI"
)

const collectionName = "httpSessions"

func Init(db *mgo.Database, twitchAPIClient *twitchAPI.Client) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: 15 * time.Minute})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"key"}})

	return &Service{collection: collection, twitchAPIClient: twitchAPIClient}
}
