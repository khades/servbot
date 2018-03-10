package httpbackend

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/khades/servbot/eventbus"
	"github.com/sirupsen/logrus"
)

func websocketEventbusWriter(w http.ResponseWriter, r *http.Request, messageID string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "websocket",
		"action":  messageID})
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }}

	conn, err := upgrader.Upgrade(w, r, nil)
	pongWait := 40 * time.Second

	if err != nil {
		logger.Debug(err)
		return
	}

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ping := func(value string) {
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
	eventbus.EventBus.Subscribe("ping", ping)
	defer eventbus.EventBus.Unsubscribe("ping", ping)

	write := func(value string) {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(value)); err != nil {
			return
		}
	}

	eventbus.EventBus.Subscribe(messageID, write)
	defer eventbus.EventBus.Unsubscribe(messageID, write)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
	}
}
