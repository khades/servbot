package repos

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/khades/servbot/httpclient"
	cache "github.com/patrickmn/go-cache"
)

var userIDCacheObject = cache.New(60*time.Minute, 600*time.Second)

type twitchUserInfo struct {
	ID          string `json:"_id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

type usersLoginStruct struct {
	Users []twitchUserInfo `json:"users"`
}

func GetUsernameByID(userID *string) (*string, error) {
	value, found := userIDCacheObject.Get("id-" + *userID)
	if found == false {
		usersString := "https://api.twitch.tv/kraken/users/" + *userID
		resp, error := httpclient.TwitchV5(*userID, "GET", usersString, nil)
		if error != nil {
			return nil, error
		}
		defer resp.Body.Close()
		var twitchUser twitchUserInfo
		marshallError := json.NewDecoder(resp.Body).Decode(&twitchUser)
		if marshallError != nil {
			return nil, marshallError
		}
		if twitchUser.Name == "" {
			return nil, errors.New("Not found")
		}

		userName := strings.ToLower(twitchUser.Name)
		userIDCacheObject.Set("username-"+userName, *userID, 0)
		userIDCacheObject.Set("id-"+*userID, userName, 0)
		return &userName, nil

	} else {
		result := value.(string)
		return &result, nil
	}
}
func GetUsersID(users *[]string) (*map[string]string, error) {
	notFoundUsers := []string{}
	result := make(map[string]string)
	log.Println(*users)
	log.Printf("Users: %d", len(*users))
	for _, user := range *users {
		value, found := userIDCacheObject.Get("username-" + user)
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
		log.Printf("Found users: %d", len(result))
		return &result, nil
	}
	sliceStart := 0
	log.Println(notFoundUsers)
	for index, _ := range notFoundUsers {
		if index == len(notFoundUsers)-1 || index-sliceStart > 48 {
			log.Printf("%d - %d", sliceStart, index)
			usersString := "https://api.twitch.tv/kraken/users?login=" + strings.Join(notFoundUsers[sliceStart:index+1], ",")
			resp, error := httpclient.TwitchV5(Config.ClientID, "GET", usersString, nil)
			if error != nil {
				return nil, error
			}
			defer resp.Body.Close()
			var usersWithID usersLoginStruct
			marshallError := json.NewDecoder(resp.Body).Decode(&usersWithID)
			if marshallError != nil {
				return nil, marshallError
			}
			log.Printf("That request returned %d users", len(usersWithID.Users))
			for _, user := range usersWithID.Users {
				result[user.DisplayName] = user.ID
				userIDCacheObject.Set("username-"+strings.ToLower(user.Name), user.ID, 0)
				userIDCacheObject.Set("id-"+user.ID, strings.ToLower(user.Name), 0)
			}
			if len(usersWithID.Users) == 0 {
				log.Println(len(notFoundUsers[sliceStart : index+1]))
				log.Println(sliceStart)
				log.Println(index + 1)
				for _, user := range notFoundUsers[sliceStart : index+1] {
					userIDCacheObject.Set("username-"+user, "rejected", 600*time.Minute)
				}
			}
			sliceStart = index

		}
	}
	// if len(result) != len(*users) {
	// 	for _, user := range *users {
	// 		if result[user] == "" {
	// 		}
	// 	}
	// }
	log.Printf("Found users: %d", len(result))
	return &result, nil
}
