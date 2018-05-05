package repos

import (
	"sort"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

var usernameCacheCollection = "usernameCacheColleciton"

type usernameCache struct {
	UserID    string
	User      string
	CreatedAt time.Time
}

var usernameCacheChatDates = make(map[string]time.Time)
var usernameCacheRejectDates = make(map[string]time.Time)

func updateUserToUserIDFromChat(userID *string, user *string) {
	date, found := usernameCacheChatDates[*userID]
	if found == false || time.Now().Sub(date) < 60*time.Minute {
		now := time.Now()
		usernameCacheChatDates[*userID] = now
		updateUserToUserID(userID, user, now)
	}
}

func getUsersByIDFromDB(users []string) ([]usernameCache, error) {
	var result = []usernameCache{}
	error := db.C(usernameCacheCollection).Find(bson.M{"user": bson.M{"$in": users}}).All(&result)
	return result, error
}

func getUsersByUsernameFromDB(userIDs []string) ([]usernameCache, error) {
	var result = []usernameCache{}
	error := db.C(usernameCacheCollection).Find(bson.M{"userid": bson.M{"$in": userIDs}}).All(&result)
	return result, error
}

func updateUserToUserID(userID *string, user *string, createdAt time.Time) {
	db.C(usernameCacheCollection).Upsert(bson.M{"userid": *userID}, bson.M{"$set": usernameCache{*userID, *user, createdAt}})
}

// GetChannelNameByID tries to fetch channelname for specified channelid from cache, falling back to channelInfo database
func GetChannelNameByID(channelID *string) (*string, error) {

	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "GetUsersID",
		"action":  "GetChannelNameByID"})
	logger.Debugf("Function is called")
	logger.Debugf("Looking for channel by id %s", *channelID)

	channel, channelError := GetChannelInfo(channelID)

	if channelError != nil {
		logger.Debugf("Looking for channelid %s in database failed: %s", *channelID, channelError.Error())

		return nil, channelError
	}

	userName := strings.ToLower(channel.Channel)
	logger.Debugf("ChannelID %s has name: %s", *channelID, strings.ToLower(channel.Channel))

	return &userName, nil
}

// GetUsersID fetches usernames from twitch for specified users, and caches them for 6 hours, resolving username->id
func GetUsersID(users []string) (*map[string]string, error) {
	// Cache every username -> userid pair for month
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "GetUsersID",
		"action":  "GetUsersID"})
	logger.Debugf("Function is called")

	result := make(map[string]string)
	logger.Debugf("Input users length: %d", len(users))
	logger.Debugf("Users: %s", strings.Join(users, ", "))

	for user, userRejectionDate := range usernameCacheRejectDates {
		if time.Now().Sub(userRejectionDate) < 20*time.Minute {
			index := sort.SearchStrings(users, user)
			if index < len(users) {
				users = remove(users, index)
			}
		}
	}
	logger.Debugf("Users after rejection: %s", strings.Join(users, ", "))

	usernamesDB, error := getUsersByUsernameFromDB(users)

	if error == nil {
		logger.Debugf("Users found: %+v", usernamesDB)

		for _, user := range usernamesDB {
			//	if time.Now().Sub(user.CreatedAt) < 60*time.Minute {
			// logger.Debugf("User found in database: %s", user.User)

			result[user.User] = user.UserID
			index := sort.SearchStrings(users, user.User)
			if index < len(users) {
				users = remove(users, index)
			}
			//	}
		}
	}

	if len(users) == 0 {
		logger.Debugf("All users found")
		logger.Debugf("Result: %+v", result)
		return &result, nil
	}

	logger.Debugf("Total found users: %d", len(result))
	logger.Debugf("Total not found users: %d", len(users))

	twitchUsers, usersError := getTwitchUsersByDisplayName(users)

	if usersError != nil {
		logger.Debugf("External function error encountered: %s", usersError.Error())
		return nil, usersError

	}

	for _, user := range twitchUsers {
		username := strings.ToLower(user.DisplayName)
		logger.Debugf("Found user %s with id %s", user.DisplayName, user.ID)

		result[username] = user.ID
		updateUserToUserID(&user.ID, &username, time.Now())
		index := sort.SearchStrings(users, username)
		if index < len(users) {
			users = remove(users, index)
		}
	}

	for _, user := range users {
		usernameCacheRejectDates[user] = time.Now()
	}

	logger.Debugf("Not found users: %s", strings.Join(users, ", "))

	logger.Debugf("Returning %d users", len(result))
	logger.Debugf("Result: %+v", result)

	return &result, nil
}

// GetUsernames fetches userID from twitch for specified users, and caches them for 6 hours
func GetUsernames(userIDs []string) (*map[string]string, error) {
	// Cache every username -> userid pair for month
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "GetUsersID",
		"action":  "GetUsernames"})
	logger.Debugf("Function is called")

	result := make(map[string]string)
	logger.Debugf("Input users length: %d", len(userIDs))
	logger.Debugf("Users: %s", strings.Join(userIDs, ", "))

	usernamesDB, error := getUsersByIDFromDB(userIDs)

	if error == nil {
		logger.Debugf("Users found: %+v", usernamesDB)

		for _, user := range usernamesDB {
			//	if time.Now().Sub(user.CreatedAt) < 60*time.Minute {
			// logger.Debugf("User found in database: %s", user.User)

			result[user.UserID] = user.User
			index := sort.SearchStrings(userIDs, user.UserID)
			if index < len(userIDs) {
				userIDs = remove(userIDs, index)
			}
			//	}
		}
	}

	if len(userIDs) == 0 {
		logger.Debugf("All users found")
		logger.Debugf("Result: %+v", result)

		return &result, nil
	}

	logger.Debugf("Total found users: %d", len(result))
	logger.Debugf("Total not found users: %d", len(userIDs))

	twitchUsers, usersError := getTwitchUsersByID(userIDs)

	if usersError != nil {
		logger.Debugf("External function error encountered: %s", usersError.Error())
		return nil, usersError

	}

	for _, user := range twitchUsers {
		username := strings.ToLower(user.DisplayName)
		logger.Debugf("Found user %s with id %s", user.DisplayName, user.ID)

		result[user.ID] = username
		updateUserToUserID(&user.ID, &username, time.Now())
	
	}



	logger.Debugf("Not found users: %s", strings.Join(userIDs, ", "))

	logger.Debugf("Returning %d users", len(result))
	logger.Debugf("Result: %+v", result)

	return &result, nil
}
