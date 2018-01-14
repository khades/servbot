package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
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
	json.NewEncoder(w).Encode(*result)
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
	result, _ := repos.GetSubsForChannelWithLimit(channelID, date)

	json.NewEncoder(w).Encode(*result)
}
func subscriptionEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	log.Println("Staring ws")
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }}

	conn, err := upgrader.Upgrade(w, r, nil)
	pongWait := 40 * time.Second


	if err != nil {
		log.Println(err)
		return
	}

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		log.Println("Got pong")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	
	ping := func(value string) {
		log.Println(value)
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Println(err)

			return
		}
	}

	write := func(value string) {
		log.Println(value)

		if err := conn.WriteMessage(websocket.TextMessage, []byte(value)); err != nil {
			log.Println(err)
			return
		}
	}
	eventbus.EventBus.Subscribe("ping", ping)

	eventbus.EventBus.Subscribe(eventbus.EventSub(channelID), write)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
	}

	defer eventbus.EventBus.Unsubscribe(eventbus.EventSub(channelID), write)
	defer eventbus.EventBus.Unsubscribe("ping", ping)

	defer log.Println("Disconnecting Subscription Socket")

}
