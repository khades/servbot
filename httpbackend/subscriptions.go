package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"goji.io/pat"

	"github.com/JanBerktold/sse"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type subscriptionsResponse struct {
	Channel       string                    `json:"channel"`
	Subscriptions []models.SubscriptionInfo `json:"subscriptions"`
}
type subscriptionEvent struct {
	Subscription     models.SubscriptionInfo `json:"subscription"`
	CurrentCallTime  time.Time               `json:"currentCallTimetime"`
	PreviousCallTime time.Time               `json:"previousCallTimetime"`
}

func subscriptions(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	var response = subscriptionsResponse{Channel: *channelName}
	result, error := repos.GetSubsForChannel(channelID)
	if error == nil {
		response.Subscriptions = *result
	}
	json.NewEncoder(w).Encode(response)
}

func subscriptionsWithLimit(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	dateLimit := pat.Param(r, "limit")
	if dateLimit == "" {
		writeJSONError(w, "limit is not defined", http.StatusUnprocessableEntity)
		return
	}
	unixTime, error := strconv.ParseInt(dateLimit, 10, 64)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusUnprocessableEntity)
		return
	}
	//log.Println(unixTime)
	date := time.Unix(0, unixTime*int64(time.Millisecond))
	//log.Println(date)
	var response = subscriptionsResponse{Channel: *channelName}
	result, error := repos.GetSubsForChannelWithLimit(channelID, date)
	if error == nil {
		response.Subscriptions = *result
	}
	json.NewEncoder(w).Encode(response)
}
func subscriptionEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	conn, _ := sse.Upgrade(w, r)
	channel := make(chan string)
	write := func(value string) {
		channel <- value
	}
	eventbus.EventBus.On("ping "+eventbus.EventSub(channelID), write)
	for conn.IsOpen() {
		msg := <-channel
		conn.WriteString(msg)
	}
	defer eventbus.EventBus.Off("ping "+eventbus.EventSub(channelID), write)
	defer log.Println("Disconnecting Subscription SSE")

}
