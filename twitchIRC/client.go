package twitchIRC

import (
	"fmt"
	"net"
	"time"

	"github.com/khades/servbot/channelInfo"

	"github.com/khades/servbot/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

// Client struct defines object that will send messages to a twitch server
type Client struct {
	// Dependencies
	config             *config.Config
	channelInfoService *channelInfo.Service
	// Own Fields
	client          *irc.Client
	handle         func (client *Client, message *irc.Message)
	bounces         map[string]time.Time
	ready           bool
	modChannelIndex int
	messageQueue    []string
	messagesSent    int
}

func (client *Client) Start() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "bot",
		"feature": "bot",
		"action":  "Start"})
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		logger.Fatalln(err)
	}
	ircConfig := irc.ClientConfig{
		Nick:    client.config.BotUserName,
		Pass:    client.config.OauthKey,
		User:    client.config.BotUserName,
		Name:    client.config.BotUserName,
		Handler: client}

	client.client = irc.NewClient(conn, ircConfig)
	logger.Info("Bot is starting")

	clientError := client.client.Run()
	logger.Info(clientError)
	logger.Fatal("Bot died")
	client.ready = false
	conn.Close()
}

// PushMessage adds message to array of messages to prevent global bans for bot
func (client *Client) PushMessage(message string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchIRC",
		"feature": "twitchIRC",
		"action":  "PushMessage"})
	logger.Infof("Pushing message: %s", message)
	if client.messagesSent > 3 {
		client.messageQueue = append(client.messageQueue, message)
		logger.Debugf("Messages in queue :", len(client.messageQueue))

	} else {
		client.client.Write(message)
		client.messagesSent = client.messagesSent + 1
	}
}

// SendMessages gets slice of messages to send periodically, sends them and updates list of messages, need to be periodically called
func (client *Client) SendMessages(interval int) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchIRC",
		"feature": "twitchIRC",
		"action":  "SendMessages"})
	client.messagesSent = 0
	queueSliceSize := 3
	arrayLen := len(client.messageQueue)

	if arrayLen == 0 {
		return
	}
	logger.Debugf("Array length is:", arrayLen)

	if arrayLen < queueSliceSize {
		queueSliceSize = arrayLen
	}
	client.messagesSent = queueSliceSize

	messagesToSend := client.messageQueue[:queueSliceSize]
	logger.Debugf("Messages to send:", len(messagesToSend))

	for _, message := range messagesToSend {
		client.client.Write(message)
	}
	client.messageQueue = client.messageQueue[queueSliceSize:]
	if len(client.messageQueue) > 0 {
		logger.Debugf("Messaged Delayed:", len(client.messageQueue))
	}
}

// SendDebounced prevents from sending data too frequent in public chat sending it to a PM
func (client *Client) SendDebounced(message OutgoingDebouncedMessage) {
	key := fmt.Sprintf("%s-%s", message.Message.Channel, message.Command)
	if client.bounces == nil {
		client.bounces = make(map[string]time.Time)
	}
	bounce, found := client.bounces[key]
	if found && int(time.Now().Sub(bounce).Seconds()) < 15 {
		client.SendPrivate(&message.Message)
	} else {
		client.bounces[key] = time.Now()
		client.SendPublic(&OutgoingMessage{
			Channel: message.Message.Channel,
			Body:    message.Message.Body,
			User:    message.RedirectTo})
	}
}

// SendPublic writes data to a specified chat
func (client *Client) SendPublic(message *OutgoingMessage) {
	if client.ready {
		messageString := ""
		if message.User != "" {
			messageString = fmt.Sprintf("PRIVMSG #%s :@%s %s", message.Channel, message.User, message.Body)
		} else {
			messageString = fmt.Sprintf("PRIVMSG #%s :%s", message.Channel, message.Body)
		}
		client.PushMessage(messageString)
	}
}

// SendPrivate writes data in private to a user
func (client *Client) SendPrivate(message *OutgoingMessage) {
	if client.ready && message.User != "" {
		messageString := fmt.Sprintf("PRIVMSG #jtv :/w %s Channel %s: %s", message.User, message.Channel, message.Body)
		client.PushMessage(messageString)

	}
}

// SendModsCommand runs mod command, need to be periodically called
func (client *Client) SendModsCommand() {
	if len(client.config.Channels) == 0 {
		return

	}
	channelName := client.config.Channels[client.modChannelIndex]
	if channelName != "" {
		client.SendPublic(&OutgoingMessage{Channel: channelName, Body: "/mods"})
	}
	client.modChannelIndex++
	if client.modChannelIndex == len(client.config.Channels) || client.modChannelIndex > len(client.config.Channels) {
		client.modChannelIndex = 0
	}

}

func (client *Client) Handle(client *irc.Client, message *irc.Message) {
	client.handle(client, message)

	if message.Command == "001" {
		client.Write("CAP REQ twitch.tv/tags")
		client.Write("CAP REQ twitch.tv/membership")
		client.Write("CAP REQ twitch.tv/commands")
		activeChannels, _ := client.channelInfoService.GetActiveChannels()
		for _, value := range activeChannels {
			client.Write("JOIN #" + value.Channel)
		}
		client.ready = true
		client.SendModsCommand()
	}
}
