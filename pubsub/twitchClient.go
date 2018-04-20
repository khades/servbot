package pubsub

import (
	"encoding/json"

	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

var IsWorking = false

type twitchWSOutgoingMessage struct {
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

func RunPubSub() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "pubsub",
		"feature": "pubsub",
		"action":  "RunPubSub"})
	var timerDur time.Duration
	channels, channelError := repos.GetChannelsWithExtendedLogging()
	if channelError != nil && len(channels) == 0 {
		return
	}

	for {
		channels, channelError := repos.GetChannelsWithExtendedLogging()
		if channelError != nil && len(channels) == 0 {
			return
		}
		var topics []string
		for _, channel := range channels {
			topics = append(topics, "chat_moderator_actions."+repos.Config.BotUserID+"."+channel.ChannelID)
		}
		logger.Info("Starting Pubsub client")
		timer := time.NewTimer(timerDur * time.Second)
		<-timer.C
		twitchPubSubClient(topics)
		logger.Info("Pubsub client died")
		timerDur = timerDur + 5
	}
}

func twitchPubSubClient(topics []string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "pubsub",
		"feature": "pubsub",
		"action":  "twitchPubSubClient"})
	u := url.URL{Scheme: "wss", Host: "pubsub-edge.twitch.tv"}
	logger.Debugf("Connecting to %s", u.String())
	pongAfterWait := 5 * time.Minute

	netDialer := net.Dialer{
		Timeout:  5 * time.Second,
		Deadline: time.Now().Add(5 * time.Second)}

	dialer := websocket.Dialer{
		NetDial:          netDialer.Dial,
		HandshakeTimeout: 5 * time.Second,
		Proxy:            http.ProxyFromEnvironment}

	conn, _, err := dialer.Dial(u.String(), nil)

	if err != nil {
		logger.Infof("Dialing failed: %s", err.Error())
		return
	}
	pongWait := 12 * time.Second

	writeErr1 := conn.WriteJSON(twitchWSOutgoingMessage{Type: "LISTEN", Nonce: "twitchPubSub", Data: authMessageData{AuthToken: strings.Replace(repos.Config.OauthKey, "oauth:", "", 1), Topics: []string{"chat_moderator_actions." + repos.Config.BotUserID + ".40635840"}}})
	if writeErr1 != nil {
		logger.Info("Initial message error:", writeErr1.Error())
		return
	}

	// writeErr2 := conn.WriteJSON(twitchWSOutgoingMessage{Type: "PING"})
	// if writeErr2 != nil {
	// 	logger.Info("Initial ping error:", writeErr2.Error())
	// 	return
	// }

	IsWorking = true

	defer conn.Close()
	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()

	done := make(chan string)

	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				logger.Infof("Read error: %s", err.Error())
				break
			}

			logger.Debugf("Incoming message body: %s", message)

			if messageType == websocket.TextMessage {
				messageObj := wsMessage{}
				errm := json.Unmarshal(message, &messageObj)
				if errm != nil {
					logger.Infof("Read error: %s", errm.Error())
					return
				}
				if messageObj.Type == "RESPONSE" {
					err := conn.WriteJSON(twitchWSOutgoingMessage{Type: "PING"})
					if err != nil {
						logger.Infof("Ping write error: %s", err.Error())

						return
					}
					conn.SetReadDeadline(time.Now().Add(pongWait))

				}
				if messageObj.Type == "PONG" {
					conn.SetReadDeadline(time.Now().Add(pongAfterWait))
				}
				if messageObj.Type == "RECONNECT" {
					break
				}
				if strings.HasPrefix(messageObj.Data.Topic, "chat_moderator_actions") == true {
					channelID := strings.Split(messageObj.Data.Topic, ".")[2]

					moderAction := moderationActionMessage{}
					json.Unmarshal([]byte(messageObj.Data.Message), &moderAction)
					result := models.ChatMessage{
						MessageStruct: models.MessageStruct{
							Date:        time.Now(),
							Username:    moderAction.Data.Args[0],
							MessageType: strings.ToLower(moderAction.Data.ModeratorAction),
							BanIssuer:   moderAction.Data.User,
							BanIssuerID: moderAction.Data.UserID},
						User:      moderAction.Data.Args[0],
						ChannelID: channelID,
						UserID:    moderAction.Data.RecipientID}

					if moderAction.Data.ModeratorAction == "ban" {
						result.MessageStruct.BanReason = moderAction.Data.Args[1]
					}

					if moderAction.Data.ModeratorAction == "timeout" {
						length, _ := strconv.Atoi(moderAction.Data.Args[1])
						result.MessageStruct.BanLength = length
					}

					repos.LogMessage(&result)

				}
			}

		}
		logger.Infof("Websocket connection closed")

		done <- "wsended"
	}()

Loop:
	for {
		select {
		case <-done:
			break Loop
		case <-ticker.C:
			err := conn.WriteJSON(twitchWSOutgoingMessage{Type: "PING"})
			if err != nil {
				logger.Infof("Ping write error: %s", err.Error())
				break Loop
			}

			conn.SetReadDeadline(time.Now().Add(pongWait))
		}
	}
	IsWorking = false

}
