package bot

import (
	"log"

	"github.com/belak/irc"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

// IrcClient struct handles stuff we use aftr
type IrcClient struct {
	Client *irc.Client
	Ready  bool
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
	// TODO Внешний контейнер для темплейтов
	if ircClient.Ready {
		if message.User != "" {
			ircClient.Client.Write(basicTemplatesInstance.PublicTemplate.Render(message))
		} else {
			ircClient.Client.Write(basicTemplatesInstance.PublicNonTargetedTemplate.Render(message))
		}
	}
}

// SendPrivate  writes data in private to a user
func (ircClient IrcClient) SendPrivate(message models.OutgoingMessage) {
	if ircClient.Ready && message.User != "" {
		ircClient.Client.Write(basicTemplatesInstance.PrivateTemplate.Render(message))
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

// IrcClientInstance is concrete irc client we work with
var IrcClientInstance = IrcClient{Ready: false}
