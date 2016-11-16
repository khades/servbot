package bot

import (
	"log"
	"strings"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func subHandler(message *irc.Message, ircClient *ircClient.IrcClient) {
	log.Println("Got first sub")
	user := strings.Split(message.Params[1], " ")[0]
	channel := message.Params[0][1:]
	if user != "" && channel != "" {
		formedMessage := models.ChatMessage{
			Channel:          channel,
			User:             user,
			IsMod:            false,
			IsSub:            true,
			IsPrime:          strings.Contains(message.String(), "Twitch Prime"),
			Date:             time.Now(),
			SubscriptionInfo: &models.SubscriptionInfo{Count: 1}}
		repos.LogMessage(&formedMessage)
		sendSubMessage(&channel, &user)
		log.Printf("Channel %v: %v subbed\n", formedMessage.Channel, formedMessage.User)
	}
}
