package userResolve

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/twitchAPI"
)

const collectionName = "usernameCacheColleciton"

func Init(db *mgo.Database, twitchAPIClient *twitchAPI.Client) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: 60 * 12 * time.Minute})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"id"}})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"userid"}})

	service := Service{
		collection:               collection,
		twitchAPIClient:          twitchAPIClient,
		usernameCacheChatDates:   make(map[string]time.Time),
		usernameCacheRejectDates: make(map[string]time.Time),
	}
	return &service
}
