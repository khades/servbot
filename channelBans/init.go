package channelBans
import (
	"github.com/globalsign/mgo"

)

const collectionName = "channelBans"

func Init(db *mgo.Database) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection: collection,
	}
}

