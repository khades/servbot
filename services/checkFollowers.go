package services

import (
	"encoding/json"
	"strings"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type links struct {
	Next string `json:"next"`
}
type follows struct {
	User user `json:"user"`
}
type user struct {
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}
type followerResponse struct {
	Cursor  string    `json:"_cursor"`
	Links   links     `json:"_links"`
	Follows []follows `json:"follows"`
}

// CheckChannelsFollowers process followers of all channels on that instance of bot, that code will be deprecated after webhooks will be done
func CheckChannelsFollowers() error{
	channels, error := repos.GetActiveChannels()
	if error != nil {
		return error
	}
	for _, value := range channels {
		checkOneChannelFollowers(&value)
	}
	return nil
}
func checkOneChannelFollowers(channel *models.ChannelInfo) {
	cursor := ""
	cursorObject, error := repos.GetFollowerCursor(&channel.ChannelID)
	if error != nil && error.Error() != "not found" {
		return
	}

	if error != nil && error.Error() == "not found" || cursorObject.Cursor == "" {
		cursor = getInitialCursor(&channel.ChannelID)
		if cursor != "" {
			repos.SetFollowerCursor(&channel.ChannelID, &cursor)
		}

	} else {
		cursor = cursorObject.Cursor
	}
	if cursor == "" {
		return
	}

	followers, followersError := getFollowers(&channel.ChannelID, &cursor)
	if followersError != nil || followers.Cursor == "" || len(followers.Follows) == 0 {
		return
	}

	repos.SetFollowerCursor(&channel.ChannelID, &followers.Cursor)
	followersToGreet := []string{}
	for _, follow := range followers.Follows {
		alreadyGreeted, _ := repos.CheckIfFollowerGreeted(&channel.ChannelID, &follow.User.Name)
		if alreadyGreeted == false {
			followersToGreet = append(followersToGreet, follow.User.Name)

			repos.AddFollowerToList(&channel.ChannelID, &follow.User.Name)
		}

	}
	if len(followersToGreet) > 0 {

		alertInfo, alertError := repos.GetSubAlert(&channel.ChannelID)
		
	
		if  alertError == nil && alertInfo.Enabled == true && alertInfo.FollowerMessage != "" {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channel.Channel,
				Body:    "@" + strings.Join(followersToGreet, " @") + " " + alertInfo.FollowerMessage})
		}
	}

}
func getFollowers(channelID *string, cursor *string) (*followerResponse, error) {
	url := "https://api.twitch.tv/kraken/channels/" + *channelID + "/follows?direction=ASC&cursor=" + *cursor

	resp, respError := httpclient.TwitchV5(repos.Config.ClientID, "GET", url, nil)
	if respError != nil {
		return nil, respError
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var responseBody = new(followerResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)

	if marshallError != nil {
		return nil, marshallError
	}
	return responseBody, nil

}
func getInitialCursor(channelID *string) string {
	url := "https://api.twitch.tv/kraken/channels/" + *channelID + "/follows?direction=DESC&limit=1"
	resp, respError := httpclient.TwitchV5(repos.Config.ClientID, "GET", url, nil)
	if respError != nil {
		return ""
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var responseBody = new(followerResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)
	if marshallError != nil {

		return ""
	}
	return responseBody.Cursor
}
