package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func resubHandler(message *irc.Message, ircClient *ircClient.IrcClient) {
	log.Println(message.String())

	systemMsg, systemMsgFound := message.Tags.GetTag("system-msg")
	prime := systemMsgFound && strings.Contains(systemMsg, `just\ssubscribed\swith\sTwitch\sPrime`)
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
				IsPrime:   prime,
				Date:      time.Now()}

			repos.LogSubscription(&loggedSubscription)

			//channels.SubscriptionChannel <- loggedSubscription
			sendResubMessage(&channel, &channelID, &user, &resubCount)
			log.Printf("Channel %v: %v resubbed for %v months\n", channel, user, resubCount)

			eventbus.EventBus.Trigger(eventbus.EventSub(&channelID))
		}
	}
}
