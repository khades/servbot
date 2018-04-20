package httpbackend

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type twitchPubSubFollows struct {
	Data []twitchPubSubFollower `json:"data"`
}
type twitchPubSubFollower struct {
	ChannelID string    `json:"to_id"`
	UserID    string    `json:"from_id"`
	Date      time.Time `json:"followed_at"`
}

type twitchPubSubStreams struct {
	Data []twitchPubSubStream `json:"data"`
}
type twitchPubSubStream struct {
	ChannelID string    `json:"user_id"`
	GameID    string    `json:"game_id"`
	Title     string    `json:"title"`
	Viewers   int       `json:"viewer_count"`
	Start     time.Time `json:"started_at"`
}

func webhookStream(w http.ResponseWriter, r *http.Request) {
	r.Header.Get("X-Hub-Signature")
	if r.FormValue("channelID") == "" {
		log.Println("NO CHANNELID")

		return
	}
	channelID := r.FormValue("channelID")

	streams := twitchPubSubStreams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&streams)
	if err != nil {
		log.Println("PARSING ERROR")

		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	status := &models.StreamStatus{
		Online: false}

	if len(streams.Data) != 0 {
		singleStatus := streams.Data[0]
		status.GameID = singleStatus.GameID
		status.Online = true
		status.Title = singleStatus.Title
		status.Viewers = singleStatus.Viewers
		status.Start = singleStatus.Start
	}
	log.Printf("%+v", status)

	repos.PushStreamStatus(&channelID, status)
}

func webhookFollows(w http.ResponseWriter, r *http.Request) {
	followers := twitchPubSubFollows{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&followers)
	if err != nil {
		log.Println("PARSING ERROR")

		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	log.Printf("%+v", followers)
	// for _, follower := range followers.Data {
	// 	alreadyGreeted, _ := repos.CheckIfFollowerGreeted(&follower.ChannelID, &follower.UserID)
	// 	if alreadyGreeted == false {
	// 		repos.AddFollowerToList(&follower.ChannelID, &follower.UserID, follower.Date, true)
	// 	}
	// }
}

func webhookVerify(w http.ResponseWriter, r *http.Request) {
	log.Println("WE ARE THERE")

	if r.FormValue("hub.topic") == "" || r.FormValue("hub.challenge") == "" {
		io.WriteString(w, "Error")

	}
	parsedURLparts := strings.Split(strings.Replace(r.FormValue("hub.topic"), "https://api.twitch.tv/helix/", "", 1), "?")
	topic := parsedURLparts[0]
	channelID := strings.Split(parsedURLparts[1], "=")[1]
	challenge := r.FormValue("hub.challenge")
	log.Println(channelID)
	log.Println(topic)
	log.Println(challenge)
	repos.PutChallengeForWebHookTopic(&channelID, &topic, &challenge)
	io.WriteString(w, challenge)

}
