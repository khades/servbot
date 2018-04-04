package repos

import (

	//	"github.com/khades/servbot/bot"
	"sort"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
)

var followersListСollectionName = "followersList"

type followerCacheStruct struct {
	isFollower bool
	Date       time.Time
}

// GetIfFollowerToChannel is cache for checking if user followed to channel, and returns duration
func GetIfFollowerToChannel(channelID *string, userID *string) (bool, time.Time) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "followers",
		"action":  "GetIfFollowerToChannel"})

	var followerResult models.Follower
	error := db.C(followersListСollectionName).Find(bson.M{"channelid": *channelID, "userid": *userID}).One(&followerResult)

	if error == nil && followerResult.IsFollower == true {
		return true, followerResult.Date
	}
	if error == nil && followerResult.IsFollower == false && time.Now().After(followerResult.Date.Add(1* time.Hour)) {
		return false, followerResult.Date
	}
	found, date := getUserFollowDate(channelID, userID)
	logger.Debugf("Result for user %s on channel %s : %s - %s", userID, channelID, found, date.String())
	AddFollowerToList(channelID, userID, date, found)

	return found, date
}

// AddFollowerToList records user that followed channel, used to prevent follow\unfollow spam
func AddFollowerToList(channelID *string, userID *string, date time.Time, isFollower bool) {
	db.C(followersListСollectionName).Upsert(
		bson.M{"channelid": *channelID, "userid": userID},
		bson.M{"$set": bson.M{"date": date, "isfollower": isFollower}})
}

// CheckIfFollowerGreeted returns if user was already greeted on channel
func CheckIfFollowerGreeted(channelID *string, follower *string) (bool, error) {
	var result models.Follower
	error := db.C(followersListСollectionName).Find(bson.M{"channelid": *channelID, "userid": *follower}).One(&result)
	if error == nil {
		return true, nil
	}
	if error != nil && error.Error() != "not found" {
		return false, error
	}
	return false, nil
}

// CheckFollowers parses all channels for followers
func CheckFollowers() ([]models.ChannelNewFollowers, error) {
	channels, error := GetActiveChannels()
	if error != nil {
		return nil, error
	}
	result := []models.ChannelNewFollowers{}
	for _, channel := range channels {
		followers := checkOneChannelFollowers(&channel)
		if len(followers) > 0 {
			result = append(result, models.ChannelNewFollowers{
				Channel:   channel.Channel,
				ChannelID: channel.ChannelID,
				Followers: followers})
		}
	}
	return result, nil
}

func checkOneChannelFollowers(channel *models.ChannelInfo) []string {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "followers",
		"action":  "checkOneChannelFollowers"})
	cursor, _ := getFollowerCursor(&channel.ChannelID)
	logger.Debugf("Channel %s followerCursor is %s", channel.ChannelID, cursor.Cursor)
	followers, followersError := getFollowers(&channel.ChannelID, cursor.Cursor.IsZero())
	logger.Debugf("Channel %s followers repsonse: %+v", channel.ChannelID, followers)
	if followersError != nil || followers.Pagination.Cursor == "" || len(followers.Followers) == 0 {
		logger.Debugf("Channel %s has no followers", channel.ChannelID)
		return nil
	}
	sort.Sort(sort.Reverse(followers.Followers))

	setFollowerCursor(&channel.ChannelID, followers.Followers[0].Date)
	followersToGreet := []string{}
	for _, follow := range followers.Followers {
		if cursor.Cursor.IsZero() == false && cursor.Cursor.After(follow.Date) {
			continue
		}
		alreadyGreeted, _ := CheckIfFollowerGreeted(&channel.ChannelID, &follow.UserID)
		if alreadyGreeted == false {
			followersToGreet = append(followersToGreet, follow.UserID)

			AddFollowerToList(&channel.ChannelID, &follow.UserID, follow.Date, true)
		}

	}
	logger.Debugf("Channel %s followers to greet: %s", channel.ChannelID, followersToGreet)

	usersToGreet, _ := getTwitchUsersByID(followersToGreet)
	result := []string{}

	if cursor.Cursor.IsZero() {
		return result
	}

	for _, user := range usersToGreet {
		result = append(result, user.DisplayName)

	}
	return result

}
