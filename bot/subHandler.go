package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/irc.v2"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func subHandler(message *irc.Message, ircClient *ircClient.IrcClient) {
	//log.Println(message.String())
	subplanMsg, subplanMsgFound := message.Tags.GetTag("msg-param-sub-plan")
	prime := subplanMsgFound && strings.Contains(subplanMsg, "prime")
	msgID, _ := message.Tags.GetTag("msg-id")
	msgParamMonths, msgParamMonthsFound := message.Tags.GetTag("msg-param-months")
	if msgID == "sub" {
		msgParamMonths = "1"
		msgParamMonthsFound = true
	}
	user, userFound := message.Tags.GetTag("display-name")
	if userFound == false || user == "" {
		user, userFound = message.Tags.GetTag("login")
	}
	if msgID == "subgift" {
		user, userFound = message.Tags.GetTag("msg-param-recipient-display-name")
	}
	channelID, channelIDFound := message.Tags.GetTag("room-id")
	userID, userIDFound := message.Tags.GetTag("user-id")

	channel := message.Params[0][1:]
	if msgParamMonthsFound && userFound && channel != "" && channelIDFound && userIDFound {
		subCount, subCountError := strconv.Atoi(msgParamMonths)
		if subCountError == nil {
			loggedSubscription := models.SubscriptionInfo{
				User:      user,
				UserID:    userID,
				ChannelID: channelID,
				Count:     subCount,
				IsPrime:   prime,
				SubPlan:   subplanMsg,
				Date:      time.Now()}

			if subCount == 1 {
				sendSubMessage(&channel, &channelID, &user, &subplanMsg)
			} else {
				sendResubMessage(&channel, &channelID, &user, &subCount, &subplanMsg)
			}
			repos.LogSubscription(&loggedSubscription)

			log.Printf("Channel %v: %v subbed for %v months\n", channel, user, subCount)

			eventbus.EventBus.Trigger(eventbus.EventSub(&channelID))
		}
	}
}
