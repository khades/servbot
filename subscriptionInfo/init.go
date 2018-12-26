package subscriptionInfo

import (
	"github.com/globalsign/mgo"
)

const collectionName = "subscriptionInfo"

func Init(db *mgo.Database) *Service {
	collection := db.C(collectionName)
	
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})

	return &Service{
		collection: collection,
	}
}