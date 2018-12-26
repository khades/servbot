package userResolve

import (
	"sort"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/twitchAPI"
	"github.com/sirupsen/logrus"
)


type Service struct {
	collection               *mgo.Collection
	twitchAPIClient          *twitchAPI.Client
	usernameCacheChatDates   map[string]time.Time
	usernameCacheRejectDates map[string]time.Time
}



func (service *Service) Update(userID *string, user *string) {
	date, found := service.usernameCacheChatDates[*userID]
	if found == false || time.Now().Sub(date) < 60*time.Minute {
		now := time.Now()
		service.usernameCacheChatDates[*userID] = now
		service.updateUserToUserID(userID, user, now)
	}
}

func (service *Service) getUsersByIDFromDB(users []string) ([]usernameCache, error) {
	var result = []usernameCache{}
	error := service.collection.Find(bson.M{"userid": bson.M{"$in": users}}).All(&result)
	return result, error
}

func (service *Service) getUsersByUsernameFromDB(userIDs []string) ([]usernameCache, error) {
	var result = []usernameCache{}
	error := service.collection.Find(bson.M{"user": bson.M{"$in": userIDs}}).All(&result)
	return result, error
}

func (service *Service) updateUserToUserID(userID *string, user *string, createdAt time.Time) {
	service.collection.Upsert(bson.M{"userid": *userID}, usernameCache{*userID, *user, createdAt})
}

// GetUsersID fetches usernames from twitch for specified users, and caches them for 6 hours, resolving username->id
func (service *Service) GetUsersID(users []string) (*map[string]string, error) {
	// Cache every username -> userid pair for month
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "GetUsersID",
		"action":  "GetUsersID"})
	logger.Debugf("Function is called")

	result := make(map[string]string)
	logger.Debugf("Input users length: %d", len(users))
	logger.Debugf("Users: %s", strings.Join(users, ", "))

	for user, userRejectionDate := range service.usernameCacheRejectDates {
		if time.Now().Sub(userRejectionDate) < 20*time.Minute {
			index := sort.SearchStrings(users, user)
			if index < len(users) {
				users = remove(users, index)
			}
		}
	}
	logger.Debugf("Users after rejection: %s", strings.Join(users, ", "))

	usernamesDB, error := service.getUsersByUsernameFromDB(users)

	if error == nil {
		logger.Debugf("Users found: %+v", usernamesDB)

		for _, user := range usernamesDB {
			//	if time.Now().Sub(user.CreatedAt) < 60*time.Minute {
			// logger.Debugf("User found in database: %s", user.User)

			result[user.User] = user.UserID
			index := findString(users, user.User)

			if index < len(users) {

				users = remove(users, index)
			}
			//	}
		}
	} else {
		logger.Debugf("DB error: %+v", error.Error())

	}
	if len(users) == 0 {
		logger.Debugf("All users found")
		logger.Debugf("Result: %+v", result)
		return &result, nil
	}

	logger.Debugf("Total found users: %d", len(result))
	logger.Debugf("Total not found users: %d", len(users))

	twitchUsers, usersError := service.twitchAPIClient.GetTwitchUsersByDisplayName(users)

	if usersError != nil {
		logger.Debugf("External function error encountered: %s", usersError.Error())
		return nil, usersError

	}

	for _, user := range twitchUsers {
		username := strings.ToLower(user.DisplayName)
		logger.Debugf("Found user %s with id %s", username, user.ID)

		result[username] = user.ID
		service.updateUserToUserID(&user.ID, &username, time.Now())
		index := findString(users, username)
		if index < len(users) {
			users = remove(users, index)
		}
	}

	for _, user := range users {
		service.usernameCacheRejectDates[user] = time.Now()
	}

	logger.Debugf("Not found users and rejected: %s", strings.Join(users, ", "))

	logger.Debugf("Returning %d users", len(result))
	logger.Debugf("Result: %+v", result)

	return &result, nil
}

// GetUsernames fetches userID from twitch for specified users, and caches them for 6 hours
func (service *Service) GetUsernames(userIDs []string) (*map[string]string, error) {
	// Cache every username -> userid pair for month
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "GetUsersID",
		"action":  "GetUsernames"})
	logger.Debugf("Function is called")

	result := make(map[string]string)
	logger.Debugf("Input users length: %d", len(userIDs))
	logger.Debugf("Users: %s", strings.Join(userIDs, ", "))

	usernamesDB, error := service.getUsersByIDFromDB(userIDs)

	if error == nil {
		logger.Debugf("Users found: %+v", usernamesDB)

		for _, user := range usernamesDB {
			//	if time.Now().Sub(user.CreatedAt) < 60*time.Minute {
			// logger.Debugf("User found in database: %s", user.User)

			result[user.UserID] = user.User
			index := findString(userIDs, user.UserID)
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

	twitchUsers, usersError := service.twitchAPIClient.GetTwitchUsersByID(userIDs)

	if usersError != nil {
		logger.Debugf("External function error encountered: %s", usersError.Error())
		return nil, usersError

	}

	for _, user := range twitchUsers {
		username := strings.ToLower(user.DisplayName)
		logger.Debugf("Found user %s with id %s", user.DisplayName, user.ID)

		result[user.ID] = username
		service.updateUserToUserID(&user.ID, &username, time.Now())

	}

	logger.Debugf("Not found users: %s", strings.Join(userIDs, ", "))

	logger.Debugf("Returning %d users", len(result))
	logger.Debugf("Result: %+v", result)

	return &result, nil
}
