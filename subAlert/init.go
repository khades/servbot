package subAlert

import "github.com/globalsign/mgo"

const collectionName = "subAlert"

func Init(db *mgo.Database) *Service {
	collection := db.C(collectionName)
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection: collection,
	}
}
