package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var followersListСollectioName = "followersList"

func AddFollowerToList(channelID *string, follower *string) {
	Db.C(followersListСollectioName).Upsert(
		bson.M{"channelid": *channelID},
		bson.M{"$addToSet": bson.M{"followers": *follower}})
}

func CheckIfFollowerGreeted(channelID *string, follower *string) (bool, error) {

	var result models.FollowersList
	error := Db.C(followersListСollectioName).Find(bson.M{"channelid": *channelID, "followers": *follower}).One(&result)
	if error == nil {
		return true, nil
	}
	if error != nil && error.Error() != "not found" {
		return false, error
	}
	return false, nil
}
