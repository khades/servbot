package channelInfo

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/userResolve"
)

func Init(db *mgo.Database, config *config.Config, userResolveService *userResolve.Service) *Service {
	collection := db.C(channelInfoCollection)
	
	service := Service{
		collection:         collection,
		config:             config,
		userResolveService: userResolveService,
		dataArray:          make(map[string]*ChannelInfo),
	}
	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	service.PreprocessChannels()
	return &service
}
