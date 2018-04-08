package repos

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
)

var videolibraryCollection = "videolibrary"

func getVideo(videoID *string) (*models.SongRequestLibraryItem, error) {
	var result models.SongRequestLibraryItem
	err := db.C(videolibraryCollection).Find(bson.M{"videoid": *videoID}).One(&result)
	return &result, err
}
func addVideoToLibrary(videoID *string, title *string, duration *time.Duration, views int64, likes int64, dislikes int64) {
	db.C(videolibraryCollection).Upsert(bson.M{"videoid": *videoID}, bson.M{"$set": bson.M{
		"length":    *duration,
		"lastcheck": time.Now(),
		"likes":     likes,
		"views":     views,
		"dislikes":  dislikes,
		"title":     *title}})
}

func AddTagToVideo(videoID *string, tag *string) {
	db.C(videolibraryCollection).Upsert(bson.M{"videoid": *videoID}, bson.M{"$addToSet": bson.M{
		"tags": *tag}})
	db.C(songRequestCollectionName).UpdateAll(bson.M{"requests.videoid": *videoID},
		bson.M{"$addToSet": bson.M{"streamstatus.videoid.$.tags": *tag}})
}
