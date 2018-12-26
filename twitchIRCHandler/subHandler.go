package twitchIRCHandler

import (
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/subscriptionInfo"
	"strconv"
	"strings"
	"time"

	"github.com/khades/servbot/twitchIRC"

	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

func (service *TwitchIRCHandler) sub(client *twitchIRC.Client, message *irc.Message) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "sub",
		"action":  "sub"})
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
	channelInfo, _ := service.channelInfoService.Get(&channelID)
	channel := message.Params[0][1:]
	if msgParamMonthsFound && userFound && channel != "" && channelIDFound && userIDFound {
		subCount, subCountError := strconv.Atoi(msgParamMonths)
		if subCountError == nil {
			loggedSubscription := subscriptionInfo.SubscriptionInfo{
				User:      user,
				UserID:    userID,
				ChannelID: channelID,
				Count:     subCount,
				IsPrime:   prime,
				SubPlan:   subplanMsg,
				Date:      time.Now()}

			if subCount == 1 {
				service.sendSubMessage(client, channelInfo, &user, &subplanMsg)
			} else {
				service.sendResubMessage(client, channelInfo, &user, &subCount, &subplanMsg)
			}
			service.subscriptionInfoService.Log(&loggedSubscription)

			logger.Debugf("Channel %v: %v subbed for %v months", channel, user, subCount)

			service.eventBus.Publish(eventbus.EventSub(&channelID), "newsub")
			service.eventBus.Publish(eventbus.Subtrain(&channelID), "newsub")
		}
	}
}
