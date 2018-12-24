package songRequest

import (
	"github.com/asaskevich/EventBus"
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/videoLibrary"
	"github.com/khades/servbot/youtubeAPIClient"
)

var songRequestCollectionName = "songrequests"

func Init(db *mgo.Database,
youtubeAPIClient   *youtubeAPIClient.YouTubeAPIClient,
channelInfoService  *channelInfo.Service,
videoLibraryService *videoLibrary.Service,
eventBus            EventBus.Bus) *Service {
	collection := db.C(songRequestCollectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	return &Service{
		collection: collection,
		youtubeAPIClient: youtubeAPIClient,
		channelInfoService: channelInfoService,
		videoLibraryService: videoLibraryService,
		eventBus: eventBus,
	}
}
