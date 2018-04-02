package repos

import (
	"time"

	"github.com/khades/servbot/models"
	"github.com/globalsign/mgo/bson"
)

var followerCursorsCollectionName = "followerCursors"

// getFollowerCursor returns cursor of paginated followers list page, used in Twitch API requests
func getFollowerCursor(channelID *string) (*models.FollowerCursor, error) {
	var result = models.FollowerCursor{}
	error := db.C(followerCursorsCollectionName).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}

// setFollowerCursor last processed cursor of paginated followers list page, used in Twitch API requests
func setFollowerCursor(channelID *string, cursor time.Time) error {

	_, err := db.C(followerCursorsCollectionName).Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"cursor": cursor}})
	return err
}
