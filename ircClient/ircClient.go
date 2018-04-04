package ircClient

import (
	"fmt"
	"time"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

// IrcClient struct defines object that will send messages to a twitch server
type IrcClient struct {
	Client          *irc.Client
	Bounces         map[string]time.Time
	Ready           bool
	ModChannelIndex int
	MessageQueue    []string
	MessagesSent    int
}

// PushMessage adds message to array of messages to prevent global bans for bot
func (ircClient *IrcClient) PushMessage(message string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "ircclient",
		"feature": "ircClient",
		"action":  "PushMessage"})
	logger.Infof("Pushing message: %s", message)
	if ircClient.MessagesSent > 3 {
		ircClient.MessageQueue = append(ircClient.MessageQueue, message)
		logger.Debugf("Messages in queue :", len(ircClient.MessageQueue))

	} else {
		ircClient.Client.Write(message)
		ircClient.MessagesSent = ircClient.MessagesSent + 1
	}
}

// SendMessages gets slice of messages to send periodically, sends them and updates list of messages
func (ircClient *IrcClient) SendMessages(interval int) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "ircclient",
		"feature": "ircClient",
		"action":  "SendMessages"})
	ircClient.MessagesSent = 0
	queueSliceSize := 3
	arrayLen := len(ircClient.MessageQueue)

	if arrayLen == 0 {
		return
	}
	logger.Debugf("Array length is:", arrayLen)

	if arrayLen < queueSliceSize {
		queueSliceSize = arrayLen
	}
	ircClient.MessagesSent = queueSliceSize

	messagesToSend := ircClient.MessageQueue[:queueSliceSize]
	logger.Debugf("Messages to send:", len(messagesToSend))

	for _, message := range messagesToSend {
		ircClient.Client.Write(message)
	}
	ircClient.MessageQueue = ircClient.MessageQueue[queueSliceSize:]
	if len(ircClient.MessageQueue) > 0 {
		logger.Debugf("Messaged Delayed:", len(ircClient.MessageQueue))
	}
}

// SendDebounced prevents from sending data too frequent in public chat sending it to a PM
func (ircClient *IrcClient) SendDebounced(message models.OutgoingDebouncedMessage) {
	key := fmt.Sprintf("%s-%s", message.Message.Channel, message.Command)
	if ircClient.Bounces == nil {
		ircClient.Bounces = make(map[string]time.Time)
	}
	bounce, found := ircClient.Bounces[key]
	if found && int(time.Now().Sub(bounce).Seconds()) < 15 {
		ircClient.SendPrivate(&message.Message)
	} else {
		ircClient.Bounces[key] = time.Now()
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: message.Message.Channel,
			Body:    message.Message.Body,
			User:    message.RedirectTo})
	}
}

// SendRaw is wrapper to Write
func (ircClient *IrcClient) SendRaw(message string) {
	if ircClient.Ready {
		ircClient.PushMessage(message)
	}
}

// SendPublic writes data to a specified chat
func (ircClient *IrcClient) SendPublic(message *models.OutgoingMessage) {
	if ircClient.Ready {
		messageString := ""
		if message.User != "" {
			messageString = fmt.Sprintf("PRIVMSG #%s :@%s %s", message.Channel, message.User, message.Body)
		} else {
			messageString = fmt.Sprintf("PRIVMSG #%s :%s", message.Channel, message.Body)
		}
		ircClient.PushMessage(messageString)
	}
}

// SendPrivate writes data in private to a user
func (ircClient *IrcClient) SendPrivate(message *models.OutgoingMessage) {
	if ircClient.Ready && message.User != "" {
		messageString := fmt.Sprintf("PRIVMSG #jtv :/w %s Channel %s: %s", message.User, message.Channel, message.Body)
		ircClient.PushMessage(messageString)

	}
}

// SendModsCommand runs mod command
func (ircClient *IrcClient) SendModsCommand() {

	//return
	channelName := repos.Config.Channels[ircClient.ModChannelIndex]
	if channelName != "" {
		ircClient.SendPublic(&models.OutgoingMessage{Channel: channelName, Body: "/mods"})
	}
	ircClient.ModChannelIndex++
	if ircClient.ModChannelIndex == len(repos.Config.Channels) || ircClient.ModChannelIndex > len(repos.Config.Channels) {
		ircClient.ModChannelIndex = 0
	}

}
