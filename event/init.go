package event

import "github.com/globalsign/mgo"

const collectionName = "events"

func Init(db *mgo.Database) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection: collection,
	}
}
