package services

import (
	"encoding/json"
	"log"
	"net/http"
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
}

type responseBody struct {
	Stream stream
}

// CheckStreamStatus runs gettting all data from all channels bot applied to
func CheckStreamStatus() {
	for _, value := range repos.Config.Channels {
		go getStatus(value)
	}
}

func getStatus(channel string) {
	var responseBody = new(responseBody)
	url := "https://api.twitch.tv/kraken/streams/" + channel + "?client_id=c5w8oif66lg711pfrh86h8od3sek43d"
	resp, respError := http.Get(url)
	var status = models.StreamStatus{
		Online: false}
	if respError != nil {
		return
	}
	defer resp.Body.Close()

	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)
	log.Println(responseBody)
	if marshallError != nil {
		return
	}
	if !responseBody.Stream.CreatedAt.IsZero() {
		status = models.StreamStatus{
			Online:      true,
			Description: "я кот",
			Game:        responseBody.Stream.Channel.Game,
			Title:       responseBody.Stream.Channel.Status,
			Start:       responseBody.Stream.CreatedAt}
	}
	// else {
	// 	repos.SetStreamStatusOffline(channel)
	// }
	repos.PushStreamStatus(channel, status)

}

// func getOfflineStatus(channel string) {
// 	var responseBody = new(channelInfo)
// 	url := "https://api.twitch.tv/kraken/channels/" + channel
// 	resp, respError := http.Get(url)

// 	if respError != nil {
// 		return
// 	}
// 	defer resp.Body.Close()

// 	marshallError := json.NewDecoder(resp.Body).Decode(responseBody)

// 	if marshallError != nil {
// 		return
// 	}

// 	status := models.StreamStatus{
// 		Online: false,
// 		Game:   responseBody.Game,
// 		Title:  responseBody.Status}
// 	repos.PushStreamStatus(channel, status)

// }
