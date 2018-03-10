package repos

import (
	"log"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var userIDCacheObject = cache.New(360*time.Minute, 60*time.Second)

type twitchUserInfo struct {
	ID          string `json:"_id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

type usersLoginStruct struct {
	Users []twitchUserInfo `json:"users"`
}

// GetChannelNameByID tries to fetch channelname for specified channelid from cache, falling back to channelInfo database
func GetChannelNameByID(channelID *string) (*string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"action":  "GetChannelNameByID"})
	logger.Debugf("Function is called")
	logger.Debugf("Looking for channel by id %s", *channelID)
	value, found := userIDCacheObject.Get("id-" + *channelID)
	if found == true {
		result := value.(string)
		logger.Debugf("Channel by id %s is found in cache: %s", *channelID, result)

		return &result, nil
	}

	logger.Debugf("Channel by id %s is not found in cache", *channelID)

	channel, channelError := GetChannelInfo(channelID)

	if channelError != nil {
		logger.Debugf("Looking for channelid %s in database failed: %s", *channelID, channelError.Error())

		return nil, channelError
	}

	userName := strings.ToLower(channel.Channel)
	userIDCacheObject.Set("id-"+*channelID, userName, 0)
	logger.Debugf("ChannelID %s has name: %s", *channelID, strings.ToLower(channel.Channel))

	return &userName, nil

}
// GetUsersID fetches usernames from twitch for specified users, and caches them for 6 hours
func GetUsersID(users []string) (*map[string]string, error) {
	// Cache every username -> userid pair for month
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"action":  "GetUsersId"})
	logger.Debugf("Function is called")

	notFoundUsers := []string{}
	result := make(map[string]string)
	logger.Debugf("Input users length: %d", len(users))
	logger.Debugf("Users: %s", strings.Join(users, ", "))

	for _, user := range users {
		value, found := userIDCacheObject.Get("username-" + strings.ToLower(user))
		if found {
			stringValue := value.(string)
			if stringValue != "rejected" {
				result[user] = stringValue
			}
		} else {
			notFoundUsers = append(notFoundUsers, user)
		}
	}
	if len(notFoundUsers) == 0 {
		logger.Debugf("Total found users: %d, all users found", len(result))

		return &result, nil
	}
	logger.Debugf("Total found users: %d", len(result))
	logger.Debugf("Total not found users: %d", len(notFoundUsers))

	logger.Debugf("Not found users: %s", strings.Join(notFoundUsers, ", "))

	for _, user := range notFoundUsers {
		userIDCacheObject.Set("username-"+user, "rejected", 600*time.Minute)
	}

	twitchUsers, usersError := getTwitchUsersByDisplayName(notFoundUsers)
	if usersError != nil {
		logger.Debugf("External function error encountered: %s", usersError.Error())
		return nil, usersError

	}
	for _, user := range twitchUsers {
		result[strings.ToLower(user.DisplayName)] = user.ID
		userIDCacheObject.Set("username-"+strings.ToLower(user.DisplayName), user.ID, 0)
	}

	log.Printf("Returning %d users", len(result))

	return &result, nil
}
