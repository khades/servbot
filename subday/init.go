package subday

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
)

const collectionName = "subdays"

func Init(db *mgo.Database, channelInfoService *channelInfo.Service) *Service {

	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "commandname"}})

	service := &Service{
		collection:         collection,
		channelInfoService: channelInfoService,
	}


	return service
}
