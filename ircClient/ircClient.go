package ircClient

import (
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
	ircClient.SendPublic(message.Message)
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
			formedMessage, _ := basicTemplatesInstance.PublicTemplate.Render(message)

			ircClient.Client.Write(formedMessage)
		} else {
			formedMessage, _ := basicTemplatesInstance.PublicNonTargetedTemplate.Render(message)
			ircClient.Client.Write(formedMessage)
		}
	}
}

// SendPrivate  writes data in private to a user
func (ircClient IrcClient) SendPrivate(message models.OutgoingMessage) {
	if ircClient.Ready && message.User != "" {
		formedMessage, _ := basicTemplatesInstance.PrivateTemplate.Render(message)
		ircClient.Client.Write(formedMessage)
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
