package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	result, error := repos.GetSubsForChannel(channelID, time.Now())
	if error == nil {
		response.Subscriptions = *result
	}
	log.Println(error)
	json.NewEncoder(w).Encode(response)
}

func subscriptionEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	conn, _ := sse.Upgrade(w, r)
	log.Println(eventbus.EventSub(channelID))
	channel := make(chan string)
	write := func() {
		channel <- "hey"
	}
	eventbus.EventBus.On(eventbus.EventSub(channelID), write)
	for conn.IsOpen() {
		msg := <-channel
		conn.WriteString(msg)
	}
	defer eventbus.EventBus.Off(eventbus.EventSub(channelID), write)
	defer log.Println("Disconnecting Subscription SSE")

}
