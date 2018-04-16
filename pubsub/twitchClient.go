package pubsub

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/khades/servbot/repos"
)

type authMessage struct {
	Type  string          `json:"type"`
	Nonce string          `json:"nonce"`
	Data  authMessageData `json:"data"`
}
type authMessageData struct {
	AuthToken string   `json:"auth_token"`
	Topics    []string `json:"topics"`
}

type wsMessage struct {
	Type string `json:"type"`
	Data wsData `json:"data"`
}
type wsData struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
type moderationActionMessage struct {
	Data moderationActionData `json:"data"`
}
type moderationActionData struct {
	Type            string   `json:"type"`
	ModeratorAction string   `json:"moderation_action"`
	Args            []string `json:"args"`
	User            string   `json:"created_by"`
	UserID          string   `json:"created_by_user_id"`
	RecipientID     string   `json:"target_user_id"`
}

func TwitchClient() {
	u := url.URL{Scheme: "wss", Host: "pubsub-edge.twitch.tv"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	userIDs, userError := repos.GetUsersID([]string{repos.Config.BotUserName})
	if userError != nil {
		log.Fatal(userError.Error())
	}
	botID := (*userIDs)[repos.Config.BotUserName]
	conn.WriteJSON(authMessage{Type: "LISTEN", Nonce: "twitchPubSub", Data: authMessageData{AuthToken: strings.Replace(repos.Config.OauthKey, "oauth:", "", 1), Topics: []string{"chat_moderator_actions." + botID + ".40635840"}}})
	conn.WriteMessage(websocket.PingMessage, nil)
	pongWait := 11 * time.Second
	conn.SetReadDeadline(time.Now().Add(pongWait))
	defer conn.Close()
	ticker := time.NewTicker(4*time.Minute + 45*time.Second)
	defer ticker.Stop()
	pongAfterWait := 5 * time.Minute
	done := make(chan string)

	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if messageType == websocket.TextMessage {
				messageObj := wsMessage{}
				errm := json.Unmarshal(message, &messageObj)
				if errm != nil {
					log.Println("read errm:", errm)

				} else {
					if messageObj.Type == "RECONNECT" {
						break
					}
					if strings.HasPrefix(messageObj.Data.Topic, "chat_moderator_actions") == true {
						moderAction := moderationActionMessage{}
						json.Unmarshal([]byte(messageObj.Data.Message), &moderAction)
						log.Printf("%+v", moderAction)

					}
				}
				log.Println("Input is:" + string(message[:]))
			}
			if err != nil {
				log.Println("read err:", err)
				break
			}
			log.Printf("recv: %s", message)
		}
		log.Println("Ended")
		done <- "a"
	}()

	conn.SetPongHandler(func(string) error {
		log.Println("Got Pong")
		conn.SetReadDeadline(time.Now().Add(pongAfterWait))
		return nil
	})
Loop:
	for {
		select {
		case <-done:
			log.Println("Got Done")

			break Loop
		case <-ticker.C:
			log.Println("Force ping")

			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("Ping error:", err)
				break Loop
			}
			conn.SetReadDeadline(time.Now().Add(pongWait))
		}
	}
}
