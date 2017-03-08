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
	values, error := repos.GetUsersID(&[]string{channel, user})
	channelID := (*values)[channel]
	userID := (*values)[user]
	if error == nil && user != "" && userID != "" && channel != "" && channelID != "" && (strings.HasSuffix(message.String(), "subscribed!") || strings.HasSuffix(message.String(), "Twitch Prime!")) {
		loggedSubscription := models.SubscriptionInfo{
			User:      user,
			UserID:    userID,
			ChannelID: channelID,
			Count:     1,
			IsPrime:   strings.Contains(message.String(), "Twitch Prime"),
			Date:      time.Now()}
		repos.LogSubscription(&loggedSubscription)
		sendSubMessage(&channel, &channelID, &user)
		log.Printf("Channel %v: %v subbed\n", channel, user)
	}
}
