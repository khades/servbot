package httpbackend

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"goji.io/pat"

	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type subscriptionEvent struct {
	Subscription     models.SubscriptionInfo `json:"subscription"`
	CurrentCallTime  time.Time               `json:"currentCallTimetime"`
	PreviousCallTime time.Time               `json:"previousCallTimetime"`
}

func subscriptions(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	result, _ := repos.GetSubsForChannel(channelID)
	json.NewEncoder(w).Encode(result)
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
	date := time.Unix(0, unixTime*int64(time.Millisecond))
	result, _ := repos.GetSubsForChannelWithLimit(channelID, date)
	json.NewEncoder(w).Encode(result)
}

func subscriptionEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	websocketEventbusWriter(w, r, eventbus.EventSub(channelID))
}
