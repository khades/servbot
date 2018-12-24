package autoMessage

import (
	"github.com/globalsign/mgo"
)

var autoMessageCollectionName = "autoMessages"

func Init(db *mgo.Database, ) *Service {
	collection := db.C(autoMessageCollectionName)

	collection .EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	service := &Service{
		collection: collection,
	}

	return service
}
