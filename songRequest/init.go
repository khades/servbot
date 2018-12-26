package songRequest

import (
	"github.com/asaskevich/EventBus"
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/videoLibrary"
	"github.com/khades/servbot/youtubeAPI"
)

const collectionName = "songrequests"

func Init(db *mgo.Database,
	youtubeAPIClient *youtubeAPI.Client,
	channelInfoService *channelInfo.Service,
	videoLibraryService *videoLibrary.Service,
	eventBus EventBus.Bus) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection:          collection,
		youtubeAPIClient:    youtubeAPIClient,
		channelInfoService:  channelInfoService,
		videoLibraryService: videoLibraryService,
		eventBus:            eventBus,
	}
}
