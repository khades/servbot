package ircClient

import (
	"fmt"
	"log"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// IrcClient struct handles stuff we use aftr
type IrcClient struct {
	Client  *irc.Client
	Bounces map[string]time.Time
	Ready   bool
}

// SendDebounced prevents from sending data too frequent
func (ircClient IrcClient) SendDebounced(message models.OutgoingDebouncedMessage) {
	key := fmt.Sprintf("%s-%s", message.Message.Channel, message.Command)
	if ircClient.Bounces == nil {
		ircClient.Bounces = make(map[string]time.Time)
	}
	bounce, found := ircClient.Bounces[key]
	if found && int(time.Now().Sub(bounce).Seconds()) < 15 {
		ircClient.SendPrivate(message.Message)
	} else {
		ircClient.Bounces[key] = time.Now()
		ircClient.SendPublic(message.Message)
	}
}

// SendRaw is wrapper to Write
func (ircClient IrcClient) SendRaw(message string) {
	log.Println(ircClient.Ready)
	if ircClient.Ready {
		ircClient.Client.Write(message)
	}
}

// SendPublic writes data to a public chat
func (ircClient IrcClient) SendPublic(message models.OutgoingMessage) {
	if ircClient.Ready {
		if message.User != "" {
			ircClient.Client.Write(fmt.Sprintf("PRIVMSG #%s :@%s %s", message.Channel, message.User, message.Body))
		} else {
			ircClient.Client.Write(fmt.Sprintf("PRIVMSG #%s :%s", message.Channel, message.Body))
		}
	}
}

// SendPrivate  writes data in private to a user
func (ircClient IrcClient) SendPrivate(message models.OutgoingMessage) {
	if ircClient.Ready && message.User != "" {
		log.Println(message.Channel)
		ircClient.Client.Write(fmt.Sprintf("PRIVMSG #jtv :/w %s Channel #%s: %s", message.User, message.Channel, message.Body))
	}
}

// SendModsCommand runs mod command
func (ircClient IrcClient) SendModsCommand() {
	if ircClient.Ready {
		for _, value := range repos.Config.Channels {
			ircClient.SendPublic(models.OutgoingMessage{Channel: value, Body: ".mods"})
		}
	}
}
