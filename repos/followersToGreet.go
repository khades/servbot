package repos

import(
	 "github.com/globalsign/mgo/bson"

"github.com/khades/servbot/models"
)
var followersToGreetCollection = "followersToGreet"

func AddFollowerToGreetOnChannel(channelID *string, follower string) {
	db.C(followersToGreetCollection).Upsert(bson.M{"channelid": *channelID}, bson.M{"$addToSet":bson.M{"followers":follower}})
}

func ResetFollowersToGreetOnChannel(channelID *string) {
	db.C(followersToGreetCollection).Remove(bson.M{"channelid": *channelID})
}

func GetFollowersToGreet() ([]models.FollowersToGreet,error) {
	result := []models.FollowersToGreet{}
	err := db.C(followersToGreetCollection).Find(bson.M{}).All(&result)
	return result, err
}