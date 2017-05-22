package bot

import (
	"strings"

	"github.com/belak/irc"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
)

func oldSubHandler(message *irc.Message, ircClient *ircClient.IrcClient) {
	user := strings.Split(message.Params[1], " ")[0]
	channel := message.Params[0][1:]
	if strings.HasSuffix(message.String(), "subscribed!") || strings.HasSuffix(message.String(), "Twitch Prime!") {
		ircClient.SendPublic(&models.OutgoingMessage{Channel: channel,
			Body: "Не забудьте прожать кнопку над чатом, чтобы получить полноценный алерт о подписке (Если кнопки нету - обновите страницу)",
			User: user})
	}
}
