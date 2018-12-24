package channelLogs

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelBans"
	"github.com/khades/servbot/userResolve"
)

var messageLogsCollection = "messageLogs"

func Init(db *mgo.Database, channelBansService *channelBans.Service, userResolveService *userResolve.Service) *Service {
	collection := db.C(messageLogsCollection)
	db.C(messageLogsCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})

	db.C(messageLogsCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "-lastupdate"}})

	db.C(messageLogsCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "$text:knownnicknames"}})

	service:=  &Service{
		collection:         collection,
		channelBansService: channelBansService,
		userResolveService: userResolveService,
	}

	return service

}
