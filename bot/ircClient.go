package bot

import (
	"github.com/belak/irc"
	"github.com/hoisie/mustache"
	"github.com/khades/servbot/models"
)

// IrcClient struct handles stuff we use aftr
type IrcClient struct {
	Client *irc.Client
	Ready  bool
}

// SendRaw is wrapper to Write
func (ircClient IrcClient) SendRaw(message string) {
	if ircClient.Ready {
		ircClient.Client.Write(message)
	}
}

// SendPublic writes data to a public chat
func (ircClient IrcClient) SendPublic(message models.OutgoingMessage) {
	// TODO Внешний контейнер для темплейтов
	if ircClient.Ready {
		if message.User != "" {
			publicTemplate, _ := mustache.ParseString("PRIVMSG # {{ message.Channel }} @{{ message.User }} {{message.Body}")
			ircClient.Client.Write(publicTemplate.Render(message))
		} else {
			publicNonTargetedTemplate, _ := mustache.ParseString("PRIVMSG # {{ message.Channel }} {{message.Body}")
			ircClient.Client.Write(publicNonTargetedTemplate.Render(message))
		}
	}
}

// SendPrivate  writes data in private to a user
func (ircClient IrcClient) SendPrivate(message models.OutgoingMessage) {
	if ircClient.Ready && message.User != "" {
		// TODO Внешний контейнер для темплейтов
		privateTemplate, _ := mustache.ParseString("PRIVMSG #jtv /w {{ privateMessage.user }} Channel #{{ privateMessage.channe }}: {{ privateMessage.body }}")
		ircClient.Client.Write(privateTemplate.Render(message))
	}
}

// IrcClientInstance is concrete irc client we work with
var IrcClientInstance = IrcClient{Ready: false}
