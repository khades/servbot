package autoMessage

import (
	"github.com/globalsign/mgo"
)

const collectionName = "autoMessages"

func Init(db *mgo.Database) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	service := &Service{
		collection: collection,
	}

	return service
}
