package twitchIRCClient

import (
	"fmt"
	"net"
	"time"

	"github.com/khades/servbot/channelInfo"

	"github.com/khades/servbot/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

// TwitchIRCClient struct defines object that will send messages to a twitch server
type TwitchIRCClient struct {
	// Dependencies
	config             *config.Config
	channelInfoService *channelInfo.Service
	// Own Fields
	client          *irc.Client
	handle         func (client *TwitchIRCClient, message *irc.Message)
	bounces         map[string]time.Time
	ready           bool
	modChannelIndex int
	messageQueue    []string
	messagesSent    int
}

func (twitchIRCClient *TwitchIRCClient) Start() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "bot",
		"feature": "bot",
		"action":  "Start"})
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		logger.Fatalln(err)
	}
	ircConfig := irc.ClientConfig{
		Nick:    twitchIRCClient.config.BotUserName,
		Pass:    twitchIRCClient.config.OauthKey,
		User:    twitchIRCClient.config.BotUserName,
		Name:    twitchIRCClient.config.BotUserName,
		Handler: twitchIRCClient}

	twitchIRCClient.client = irc.NewClient(conn, ircConfig)
	logger.Info("Bot is starting")

	clientError := twitchIRCClient.client.Run()
	logger.Info(clientError)
	logger.Fatal("Bot died")
	twitchIRCClient.ready = false
	conn.Close()
}

// PushMessage adds message to array of messages to prevent global bans for bot
func (twitchIRCClient *TwitchIRCClient) PushMessage(message string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchIRCClient",
		"feature": "twitchIRCClient",
		"action":  "PushMessage"})
	logger.Infof("Pushing message: %s", message)
	if twitchIRCClient.messagesSent > 3 {
		twitchIRCClient.messageQueue = append(twitchIRCClient.messageQueue, message)
		logger.Debugf("Messages in queue :", len(twitchIRCClient.messageQueue))

	} else {
		twitchIRCClient.client.Write(message)
		twitchIRCClient.messagesSent = twitchIRCClient.messagesSent + 1
	}
}

// SendMessages gets slice of messages to send periodically, sends them and updates list of messages, need to be periodically called
func (twitchIRCClient *TwitchIRCClient) SendMessages(interval int) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchIRCClient",
		"feature": "twitchIRCClient",
		"action":  "SendMessages"})
	twitchIRCClient.messagesSent = 0
	queueSliceSize := 3
	arrayLen := len(twitchIRCClient.messageQueue)

	if arrayLen == 0 {
		return
	}
	logger.Debugf("Array length is:", arrayLen)

	if arrayLen < queueSliceSize {
		queueSliceSize = arrayLen
	}
	twitchIRCClient.messagesSent = queueSliceSize

	messagesToSend := twitchIRCClient.messageQueue[:queueSliceSize]
	logger.Debugf("Messages to send:", len(messagesToSend))

	for _, message := range messagesToSend {
		twitchIRCClient.client.Write(message)
	}
	twitchIRCClient.messageQueue = twitchIRCClient.messageQueue[queueSliceSize:]
	if len(twitchIRCClient.messageQueue) > 0 {
		logger.Debugf("Messaged Delayed:", len(twitchIRCClient.messageQueue))
	}
}

// SendDebounced prevents from sending data too frequent in public chat sending it to a PM
func (twitchIRCClient *TwitchIRCClient) SendDebounced(message OutgoingDebouncedMessage) {
	key := fmt.Sprintf("%s-%s", message.Message.Channel, message.Command)
	if twitchIRCClient.bounces == nil {
		twitchIRCClient.bounces = make(map[string]time.Time)
	}
	bounce, found := twitchIRCClient.bounces[key]
	if found && int(time.Now().Sub(bounce).Seconds()) < 15 {
		twitchIRCClient.SendPrivate(&message.Message)
	} else {
		twitchIRCClient.bounces[key] = time.Now()
		twitchIRCClient.SendPublic(&OutgoingMessage{
			Channel: message.Message.Channel,
			Body:    message.Message.Body,
			User:    message.RedirectTo})
	}
}

// SendPublic writes data to a specified chat
func (twitchIRCClient *TwitchIRCClient) SendPublic(message *OutgoingMessage) {
	if twitchIRCClient.ready {
		messageString := ""
		if message.User != "" {
			messageString = fmt.Sprintf("PRIVMSG #%s :@%s %s", message.Channel, message.User, message.Body)
		} else {
			messageString = fmt.Sprintf("PRIVMSG #%s :%s", message.Channel, message.Body)
		}
		twitchIRCClient.PushMessage(messageString)
	}
}

// SendPrivate writes data in private to a user
func (twitchIRCClient *TwitchIRCClient) SendPrivate(message *OutgoingMessage) {
	if twitchIRCClient.ready && message.User != "" {
		messageString := fmt.Sprintf("PRIVMSG #jtv :/w %s Channel %s: %s", message.User, message.Channel, message.Body)
		twitchIRCClient.PushMessage(messageString)

	}
}

// SendModsCommand runs mod command, need to be periodically called
func (twitchIRCClient *TwitchIRCClient) SendModsCommand() {
	if len(twitchIRCClient.config.Channels) == 0 {
		return

	}
	channelName := twitchIRCClient.config.Channels[twitchIRCClient.modChannelIndex]
	if channelName != "" {
		twitchIRCClient.SendPublic(&OutgoingMessage{Channel: channelName, Body: "/mods"})
	}
	twitchIRCClient.modChannelIndex++
	if twitchIRCClient.modChannelIndex == len(twitchIRCClient.config.Channels) || twitchIRCClient.modChannelIndex > len(twitchIRCClient.config.Channels) {
		twitchIRCClient.modChannelIndex = 0
	}

}

func (twitchIRCClient *TwitchIRCClient) Handle(client *irc.Client, message *irc.Message) {
	twitchIRCClient.handle(twitchIRCClient, message)

	if message.Command == "001" {
		client.Write("CAP REQ twitch.tv/tags")
		client.Write("CAP REQ twitch.tv/membership")
		client.Write("CAP REQ twitch.tv/commands")
		activeChannels, _ := twitchIRCClient.channelInfoService.GetActiveChannels()
		for _, value := range activeChannels {
			client.Write("JOIN #" + value.Channel)
		}
		twitchIRCClient.ready = true
		twitchIRCClient.SendModsCommand()
	}
}
