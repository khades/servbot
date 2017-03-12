package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JanBerktold/sse"
	"github.com/khades/servbot/channels"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type subscriptionsResponse struct {
	Channel       string                    `json:"channel"`
	Subscriptions []models.SubscriptionInfo `json:"subscriptions"`
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
	for {
		msg := <-channels.SubscriptionChannel
		if msg.ChannelID == *channelID {
			conn.WriteJson(msg)
		}

	}

}
