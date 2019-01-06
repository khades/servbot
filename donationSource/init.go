package donationSource

import (
	"github.com/globalsign/mgo"
)

var collectionName = "donationSources"

func Init(db *mgo.Database) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection: collection,
	}
}
