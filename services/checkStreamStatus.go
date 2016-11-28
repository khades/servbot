package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type stream struct {
	CreatedAt time.Time `json:"created_at"`
	Channel   channelInfo
}

type channelInfo struct {
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
	for _, channel := range repos.Config.Channels {
		streams[channel] = models.StreamStatus{
			Online: false}
	}

	url := "https://api.twitch.tv/kraken/streams?channel=" + strings.Join(repos.Config.Channels, ",") + "&client_id=" + repos.Config.ClientID
	log.Println(url)
	resp, respError := http.Get(url)
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
		log.Println(status)
		streams[status.Channel.Name] = models.StreamStatus{
			Online: true,
			Game:   status.Channel.Game,
			Title:  status.Channel.Status,
			Start:  status.CreatedAt}
	}
	for channel, status := range streams {
		repos.PushStreamStatus(&channel, &status)
	}
}
