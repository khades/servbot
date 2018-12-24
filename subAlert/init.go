package subAlert

import "github.com/globalsign/mgo"

var subAlertCollection = "subAlert"

func Init(db *mgo.Database) *Service {
	collection := db.C(subAlertCollection)
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection: collection,
	}
}
