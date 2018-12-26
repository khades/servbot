package followers

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/twitchAPI"

	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)


// type followerCacheStruct struct {
// 	isFollower bool
// 	Date       time.Time
// }

type Service struct {
	// Dependencies
	twitchAPIClient         *twitchAPI.Client
	followersToGreetService *followersToGreet.Service

	// Own fields
	collection *mgo.Collection
}

// IsFoller is cache for checking if user followed to channel, and returns duration
func (service *Service) IsFoller(channelID *string, userID *string) (bool, time.Time) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "followers",
		"action":  "IsFoller"})

	var followerResult Follower

	error := service.collection.Find(bson.M{"channelid": *channelID, "userid": *userID}).One(&followerResult)

	if error == nil && followerResult.IsFollower == true {
		return true, followerResult.Date
	}
	
	if error == nil && followerResult.IsFollower == false && time.Now().After(followerResult.Date.Add(1* time.Hour)) {
		return false, followerResult.Date
	}

	found, date := service.twitchAPIClient.GetUserFollowDate(channelID, userID)

	logger.Debugf("Result for user %s on channel %s : %s - %s", &userID, &channelID, found, date.String())
	service.Add(channelID, userID, date, found, false)

	return found, date
}

// Add records user that followed channel, used to prevent follow\unfollow spam
func (service *Service) Add(channelID *string, userID *string, date time.Time, isFollower bool, notify bool) {
	changeInfo, changeError := service.collection.Upsert(
		bson.M{"channelid": *channelID, "userid": userID},
		bson.M{"$set": bson.M{"date": date, "isfollower": isFollower}})
	if notify == true && changeError == nil && changeInfo.Updated == 0 {
		service.followersToGreetService.Add(channelID, *userID)
	}
}

// IsGreeted returns if user was already greeted on channel
func (service *Service) IsGreeted(channelID *string, follower *string) (bool, error) {
	var result Follower
	error := service.collection.Find(bson.M{"channelid": *channelID, "userid": *follower}).One(&result)
	if error == nil {
		return true, nil
	}
	if error != nil && error.Error() != "not found" {
		return false, error
	}
	return false, nil
}

// TODO: That is needed to be in followerAlertService
// CheckFollowers parses all channels for followers
// func (service *Service) CheckFollowers() ([]models.ChannelNewFollowers, error) {
// 	channels, error := GetActiveChannels()
// 	if error != nil {
// 		return nil, error
// 	}
// 	result := []models.ChannelNewFollowers{}
// 	for _, channel := range channels {
// 		followers := checkOneChannelFollowers(&channel)
// 		if len(followers) > 0 {
// 			result = append(result, models.ChannelNewFollowers{
// 				Channel:   channel.Channel,
// 				ChannelID: channel.ChannelID,
// 				Followers: followers})
// 		}
// 	}
// 	return result, nil
// }

// // TODO: That is needed to be in followerAlertService
// func (service *Service) checkOneChannelFollowers(channel *models.ChannelInfo) []string {
// 	logger := logrus.WithFields(logrus.Fields{
// 		"package": "repos",
// 		"feature": "followers",
// 		"action":  "checkOneChannelFollowers"})
// 	cursor, _ := getFollowerCursor(&channel.ChannelID)
// 	logger.Debugf("Channel %s followerCursor is %s", channel.ChannelID, cursor.Cursor)
// 	followers, followersError := getFollowers(&channel.ChannelID, cursor.Cursor.IsZero())
// 	logger.Debugf("Channel %s followers repsonse: %+v", channel.ChannelID, followers)
// 	if followersError != nil || followers.Pagination.Cursor == "" || len(followers.Followers) == 0 {
// 		logger.Debugf("Channel %s has no followers", channel.ChannelID)
// 		return nil
// 	}
// 	sort.Sort(sort.Reverse(followers.Followers))

// 	setFollowerCursor(&channel.ChannelID, followers.Followers[0].Date)
// 	followersToGreet := []string{}
// 	for _, follow := range followers.Followers {
// 		if cursor.Cursor.IsZero() == false && cursor.Cursor.After(follow.Date) {
// 			continue
// 		}
// 		alreadyGreeted, _ := IsGreeted(&channel.ChannelID, &follow.UserID)
// 		if alreadyGreeted == false {
// 			followersToGreet = append(followersToGreet, follow.UserID)

// 			Add(&channel.ChannelID, &follow.UserID, follow.Date, true)
// 		}

// 	}
// 	logger.Debugf("Channel %s followers to greet: %s", channel.ChannelID, followersToGreet)

// 	usersToGreet, _ := getTwitchUsersByID(followersToGreet)
// 	result := []string{}

// 	if cursor.Cursor.IsZero() {
// 		return result
// 	}

// 	for _, user := range usersToGreet {
// 		result = append(result, user.DisplayName)

// 	}
// 	return result

// }
