package httpbackend

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

type twitchPubSubFollows struct {
	Data twitchPubSubFollower `json:"data"`
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
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpbackend",
		"feature": "webhook",
		"action":  "webhookStream"})
	logger.Debugf("Request signature is %s", r.Header.Get("X-Hub-Signature"))
	if r.FormValue("channelID") == "" {
		logger.Debugf("No channel set")
		return
	}
	channelID := r.FormValue("channelID")

	streams := twitchPubSubStreams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&streams)
	if err != nil {
		logger.Debugf("JSON decode error: %s", err.Error())

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

	logger.Debugf("New status for channel %s : %+v", channelID, status)

	repos.PushStreamStatus(&channelID, status)
}

func webhookFollows(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpbackend",
		"feature": "webhook",
		"action":  "webhookFollows"})
	follower := twitchPubSubFollows{}
	logger.Debugf("Request signature is %s", r.Header.Get("X-Hub-Signature"))
	dump, dumpErr := httputil.DumpRequest(r, true)
	if dumpErr == nil {
		logger.Debugf("Repsonse is %q", dump)
	}
	decoder := json.NewDecoder(r.Body)
	// topic := "follows"
	// topicItem, topicError := repos.GetWebHookTopic(&channelID, &topic )
	err := decoder.Decode(&follower)
	if err != nil {
		logger.Debugf("JSON decode error: %s", err.Error())

		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	logger.Debugf("Incoming followers: %+v", follower)

	logger.Debugf("User %s follows channel %s", follower.Data.UserID, follower.Data.ChannelID)

	alreadyGreeted, _ := repos.CheckIfFollowerGreeted(&follower.Data.ChannelID, &follower.Data.UserID)
	if alreadyGreeted == false {
		repos.AddFollowerToList(&follower.Data.ChannelID, &follower.Data.UserID, follower.Data.Date, true)
	}

}

func webhookVerify(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("hub.topic") == "" || r.FormValue("hub.challenge") == "" {
		io.WriteString(w, "Error")
	}

	challenge := r.FormValue("hub.challenge")
	// parsedURLparts := strings.Split(strings.Replace(r.FormValue("hub.topic"), "https://api.twitch.tv/helix/", "", 1), "?")
	// topic := parsedURLparts[0]
	// channelID := strings.Split(parsedURLparts[1], "=")[1]
	// repos.PutChallengeForWebHookTopic(&channelID, &topic, &challenge)
	io.WriteString(w, challenge)
}
