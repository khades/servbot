package followersToGreet

import "github.com/globalsign/mgo"

var followersToGreetCollection = "followersToGreet"

func Init(db *mgo.Database) *Service {
	collection:= db.C(followersToGreetCollection )

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{collection: collection}
}