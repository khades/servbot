package bot

import (
	"strconv"
	"strings"
	"time"

	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

func subHandler(message *irc.Message, ircClient *ircClient.IrcClient) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "subHandler",
		"action":  "subHandler"})
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
	channelInfo, _ := repos.GetChannelInfo(&channelID)
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
				sendSubMessage(channelInfo, &user, &subplanMsg)
			} else {
				sendResubMessage(channelInfo, &user, &subCount, &subplanMsg)
			}
			repos.LogSubscription(&loggedSubscription)

			logger.Debugf("Channel %v: %v subbed for %v months", channel, user, subCount)

			eventbus.EventBus.Publish(eventbus.EventSub(&channelID))
			eventbus.EventBus.Publish(eventbus.Subtrain(&channelID), "newsub")
		}
	}
}
