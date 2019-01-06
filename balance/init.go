package balance

import "github.com/globalsign/mgo"

const collectionName = "balance"

func Init(db *mgo.Database) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})
	
	return &Service{collection}
}
