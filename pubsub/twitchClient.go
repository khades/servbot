package pubsub

import (
	"encoding/json"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/config"

	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	IsWorking          bool
	channelInfoService *channelInfo.Service
	config             *config.Config
	channelLogsService *channelLogs.Service
}

// TODO: Needs autorestart
func (client *Client) run() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "pubsub",
		"action":  "Run"})
	var timerDur time.Duration = 10 * time.Second
	var timerIncrement = func() {
		if timerDur < 1*time.Minute {
			timerDur = timerDur + (5 * time.Second)
		}
	}
	//channels, channelError := client.channelInfoService.GetChannelsWithExtendedLogging()
	//if channelError != nil && len(channels) == 0 {
	//	return
	//}

	for {
		timer := time.NewTimer(timerDur)

		channels, channelError := client.channelInfoService.GetChannelsWithExtendedLogging()
		if channelError != nil || len(channels) == 0 {
			logger.Info("Nothing to listen")
			<-timer.C
			continue
		}
		var topics []string
		for _, channel := range channels {
			topics = append(topics, "chat_moderator_actions."+client.config.BotUserID+"."+channel.ChannelID)
		}
		logger.Infof("Starting Pubsub client with topics: %+v", topics)
		<-timer.C
		client.twitchPubSubClient(topics)
		logger.Info("Pubsub client died")
		timerIncrement()
	}
}

func (client *Client) twitchPubSubClient(topics []string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "pubsub",
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

	writeErr1 := conn.WriteJSON(twitchWSOutgoingMessage{Type: "LISTEN", Nonce: "twitchPubSub", Data: authMessageData{AuthToken: strings.Replace(client.config.OauthKey, "oauth:", "", 1), Topics: topics}})
	if writeErr1 != nil {
		logger.Info("Initial message error:", writeErr1.Error())
		return
	}

	client.IsWorking = true

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
					if len(moderAction.Data.Args) == 0 {
						logger.Infof("Got unknown action, but cant get its data: %+v", messageObj)
						logger.Infof("Got unknown action, but cant get its data: %+v", moderAction)

					} else {
						result := chatMessage.ChatMessage{
							MessageStruct: chatMessage.MessageStruct{
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
							client.channelLogsService.Log(&result)

						}

						if moderAction.Data.ModeratorAction == "timeout" {
							length, _ := strconv.Atoi(moderAction.Data.Args[1])
							result.MessageStruct.BanLength = length
							client.channelLogsService.Log(&result)
						}
					}
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
	client.IsWorking = false

}
