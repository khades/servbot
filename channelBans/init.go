package channelBans
import (
	"github.com/globalsign/mgo"

)

const channelBansCollectionName = "channelBans"

func Init(db *mgo.Database) *Service {
	collection := db.C(channelBansCollectionName)

	db.C(channelBansCollectionName).EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection: collection,
	}
}

