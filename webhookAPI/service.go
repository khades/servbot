package webhookAPI

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/followers"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/streamStatus"
	"github.com/khades/servbot/webhook"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/sirupsen/logrus"
)

type Service struct {
	webHookService *webhook.Service
	streamStatusService *streamStatus.Service
	followersService *followers.Service

}
func (service *Service) streamStatus(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpbackend",
		"feature": "webhook",
		"action":  "streamStatus"})
	logger.Debugf("Request signature is %s", r.Header.Get("X-Hub-Signature"))
	dump, dumpErr := httputil.DumpRequest(r, true)
	if dumpErr == nil {
		logger.Debugf("Repsonse is %q", dump)
	}
	if r.FormValue("channelID") == "" {
		logger.Debugf("No channel set")
		return
	}
	channelID := r.FormValue("channelID")
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	topicItem, topicError := service.webHookService.Get(&channelID, "streams")
	if topicError != nil {
		logger.Debugf("Topic doesnt exists, exiting")
		httpAPI.WriteJSONError(w, "Wrong signature", http.StatusUnprocessableEntity)
		return
	}

	mac := hmac.New(sha256.New, []byte(topicItem.Secret))
	mac.Write(bodyBytes)

	logger.Debugf("calculated signature is %s", hex.EncodeToString(mac.Sum(nil)))
	logger.Debugf("Request signature is %s", r.Header.Get("X-Hub-Signature"))

	if strings.Replace(r.Header.Get("X-Hub-Signature"), "sha256=", "", 1) != hex.EncodeToString(mac.Sum(nil)) {
		logger.Debugf("Hexes are not equal, exiting")
		httpAPI.WriteJSONError(w, "Wrong signature", http.StatusUnprocessableEntity)
		return
	}
	streams := twitchPubSubStreams{}
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	err := decoder.Decode(&streams)
	if err != nil {
		logger.Debugf("JSON decode error: %s", err.Error())

		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	status := &channelInfo.StreamStatus{
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

	service.streamStatusService.Push(&channelID, status)
}

func (service *Service) follows(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpbackend",
		"feature": "webhook",
		"action":  "follows"})
	follower := twitchPubSubFollows{}
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	// dump, dumpErr := httputil.DumpRequest(r, true)
	// if dumpErr == nil {
	// 	logger.Debugf("Repsonse is %q", dump)
	// }
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))

	err := decoder.Decode(&follower)
	if err != nil {
		logger.Debugf("JSON decode error: %s", err.Error())

		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	topicItem, topicError := service.webHookService.Get(&follower.Data.ChannelID, "follows")
	if topicError != nil {
		logger.Debugf("Topic doesnt exists, exiting")
		httpAPI.WriteJSONError(w, "Wrong signature", http.StatusUnprocessableEntity)
		return
	}

	mac := hmac.New(sha256.New, []byte(topicItem.Secret))
	mac.Write(bodyBytes)

	logger.Debugf("calculated signature is %s", hex.EncodeToString(mac.Sum(nil)))
	logger.Debugf("Request signature is %s", r.Header.Get("X-Hub-Signature"))

	if strings.Replace(r.Header.Get("X-Hub-Signature"), "sha256=", "", 1) != hex.EncodeToString(mac.Sum(nil)) {
		logger.Debugf("Hexes are not equal, exiting")
		httpAPI.WriteJSONError(w, "Wrong signature", http.StatusUnprocessableEntity)
		return
	}
	logger.Debugf("Hexes are equal, proceeding")

	logger.Debugf("User %s follows channel %s", follower.Data.UserID, follower.Data.ChannelID)


		service.followersService.Add(&follower.Data.ChannelID, &follower.Data.UserID, follower.Timestamp, true, true)



}

func (service *Service) verify(w http.ResponseWriter, r *http.Request) {
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
