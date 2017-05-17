package repos

import (
	"log"

	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var followerCursorsCollectionName = "followerCursors"

func GetFollowerCursor(channelID *string) (*models.FollowerCursor, error) {
	var result = models.FollowerCursor{}
	error := Db.C(followerCursorsCollectionName).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}

func SetFollowerCursor(channelID *string, cursor *string) {

	_, err := Db.C(followerCursorsCollectionName).Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"cursor": *cursor}})
	log.Println(err)
}
