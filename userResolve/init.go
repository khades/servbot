package userResolve

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/twitchAPIClient"
)

var usernameCacheCollection = "usernameCacheColleciton"

func Init(db *mgo.Database, twitchAPIClient *twitchAPIClient.TwitchAPIClient) *Service {
	
	db.C(usernameCacheCollection).EnsureIndex(mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: 60 * 12 * time.Minute})

	db.C(usernameCacheCollection).EnsureIndex(mgo.Index{
		Key: []string{"id"}})

	db.C(usernameCacheCollection).EnsureIndex(mgo.Index{
		Key: []string{"userid"}})

	service := Service{
		collection:               db.C(usernameCacheCollection),
		twitchAPIClient:                twitchAPIClient,
		usernameCacheChatDates:   make(map[string]time.Time),
		usernameCacheRejectDates: make(map[string]time.Time),
	}
	return &service
}
