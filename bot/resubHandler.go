package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func resubHandler(message *irc.Message, ircClient *ircClient.IrcClient) {
	msgParamMonths, msgParamMonthsFound := message.Tags.GetTag("msg-param-months")
	user, userFound := message.Tags.GetTag("display-name")
	channelID, channelIDFound := message.Tags.GetTag("room-id")
	userID, userIDFound := message.Tags.GetTag("user-id")

	channel := message.Params[0][1:]
	if msgParamMonthsFound && userFound && channel != "" && channelIDFound && userIDFound {
		resubCount, resubCountError := strconv.Atoi(msgParamMonths)
		if resubCountError == nil {
			loggedSubscription := models.SubscriptionInfo{
				User:      user,
				UserID:    userID,
				ChannelID: channelID,
				Count:     resubCount,
				IsPrime:   strings.Contains(message.String(), "Twitch Prime"),
				Date:      time.Now()}

			repos.LogSubscription(&loggedSubscription)
			//channels.SubscriptionChannel <- loggedSubscription
			sendResubMessage(&channel, &channelID, &user, &resubCount)
			log.Printf("Channel %v: %v resubbed for %v months\n", channel, user, resubCount)
		}
	}
}
