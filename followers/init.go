package followers

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/twitchAPI"
)

const collectionName  = "followersList"

func Init(db *mgo.Database, twitchAPIClient *twitchAPI.Client,followersToGreetService *followersToGreet.Service) *Service {
	collection:= db.C(collectionName)
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})
	return &Service{collection: collection, twitchAPIClient: twitchAPIClient, followersToGreetService: followersToGreetService}
}

