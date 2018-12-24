package followers

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/twitchAPIClient"
)

var followersListСollectionName = "followersList"

func Init(db *mgo.Database, twitchAPIClient *twitchAPIClient.TwitchAPIClient,followersToGreetService *followersToGreet.Service) *Service {
	collection:= db.C(followersListСollectionName)
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})
	return &Service{collection: collection, twitchAPIClient: twitchAPIClient, followersToGreetService: followersToGreetService}
}

