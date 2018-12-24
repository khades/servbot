package subscriptionInfo

import (
	"github.com/globalsign/mgo"
)

var subscriptionInfoCollection = "subscriptionInfo"

func Init(db *mgo.Database) *Service {
	collection := db.C(subscriptionInfoCollection)
	
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})

	return &Service{
		collection: collection,
	}
}