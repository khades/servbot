package services

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type stream struct {
	CreatedAt time.Time `json:"created_at"`
	Channel   channelInfo
}

type channelInfo struct {
	ID     int `json:"_id"`
	Status string
	Game   string
	Name   string
}

type responseBodyStruct struct {
	Streams []stream
}

// CheckStreamStatus runs gettting all data from all channels bot applied to
func CheckStreamStatus() {
	streams := make(map[string]models.StreamStatus)
	users, error := repos.GetUsersID(&repos.Config.Channels)
	if error != nil {
		return
	}
	userIDs := []string{}
	for _, id := range *users {
		userIDs = append(userIDs, id)
	}
	for _, channel := range userIDs {
		streams[channel] = models.StreamStatus{
			Online: false}

	}

	url := "https://api.twitch.tv/kraken/streams?channel=" + strings.Join(userIDs, ",")
	resp, respError := httpclient.TwitchV5(repos.Config.ClientID, "GET", url, nil)
	if respError != nil {
		return
	}
	defer resp.Body.Close()
	var responseBody = new(responseBodyStruct)

	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)

	if marshallError != nil {
		return
	}

	for _, status := range responseBody.Streams {

		streams[strconv.Itoa(status.Channel.ID)] = models.StreamStatus{
			Online: true,
			Game:   status.Channel.Game,
			Title:  status.Channel.Status,
			Start:  status.CreatedAt}

	}
	for channel, status := range streams {
		repos.PushStreamStatus(&channel, &status)
	}
}
