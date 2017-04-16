package ircClient

import (
	"fmt"
	"log"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// IrcClient struct defines object that will send messages to a twitch server
type IrcClient struct {
	Client  *irc.Client
	Bounces map[string]time.Time
	Ready   bool
}

// SendDebounced prevents from sending data too frequent in public chat sending it to a PM
func (ircClient IrcClient) SendDebounced(message models.OutgoingDebouncedMessage) {
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
func (ircClient IrcClient) SendRaw(message string) {
	if ircClient.Ready {
		ircClient.Client.Write(message)
	}
}

// SendPublic writes data to a specified chat
func (ircClient IrcClient) SendPublic(message *models.OutgoingMessage) {
	if ircClient.Ready {
		if message.User != "" {
			ircClient.Client.Write(fmt.Sprintf("PRIVMSG #%s :@%s %s", message.Channel, message.User, message.Body))
		} else {
			ircClient.Client.Write(fmt.Sprintf("PRIVMSG #%s :%s", message.Channel, message.Body))
		}
	}
}

// SendPrivate writes data in private to a user
func (ircClient IrcClient) SendPrivate(message *models.OutgoingMessage) {
	if ircClient.Ready && message.User != "" {
		ircClient.Client.Write(fmt.Sprintf("PRIVMSG #jtv :/w %s Channel %s: %s", message.User, message.Channel, message.Body))
	}
}

// SendModsCommand runs mod command
func (ircClient IrcClient) SendModsCommand() {
	log.Println("Sending MODS")
	if ircClient.Ready {
		for _, value := range repos.Config.Channels {
			ircClient.SendPublic(&models.OutgoingMessage{Channel: value, Body: "/mods"})
		}
	}
}
