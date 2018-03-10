package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var followerCursorsCollectionName = "followerCursors"

// GetFollowerCursor returns cursor of paginated followers list page, used in Twitch API requests
func GetFollowerCursor(channelID *string) (*models.FollowerCursor, error) {
	var result = models.FollowerCursor{}
	error := db.C(followerCursorsCollectionName).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}

// SetFollowerCursor last processed cursor of paginated followers list page, used in Twitch API requests
func SetFollowerCursor(channelID *string, cursor *string) error{

	_, err := db.C(followerCursorsCollectionName).Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"cursor": *cursor}})
	return err
}
