package channelLogs

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelBans"
	"github.com/khades/servbot/userResolve"
)

const collectionName = "messageLogs"

func Init(db *mgo.Database, channelBansService *channelBans.Service, userResolveService *userResolve.Service) *Service {
	collection := db.C(collectionName)
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "-lastupdate"}})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "$text:knownnicknames"}})

	service:=  &Service{
		collection:         collection,
		channelBansService: channelBansService,
		userResolveService: userResolveService,
	}

	return service

}
