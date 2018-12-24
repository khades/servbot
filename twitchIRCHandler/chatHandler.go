package twitchIRCHandler

import (
	"github.com/khades/servbot/subscriptionInfo"
	"strconv"
	"strings"
	"time"

	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/subday"
	"github.com/khades/servbot/twitchIRCClient"
	"github.com/khades/servbot/userResolve"

	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/pubsub"
	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

type TwitchIRCHandler struct {
	subdayService      *subday.Service
	channelInfoService *channelInfo.Service
	subAlertService    *subAlert.Service
	channelLogsService *channelLogs.Service
	autoMessageService *autoMessage.Service
	userResolveService *userResolve.Service
	subscriptionInfoService *subscriptionInfo.Service
}

func (service *TwitchIRCHandler) Handle(client *twitchIRCClient.TwitchIRCClient, message *irc.Message) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "bot",
		"feature": "bot",
		"action":  "chatHandler"})
	if strings.Contains(message.String(), ":jtv") {
		logger.Debug("JTV: " + message.String())
	}
	if strings.Contains(message.String(), ":tmi.twitch.tv") {
		logger.Debug("TMI.TWITCH.TV: " + message.String())
	}
	msgID, found := message.Tags.GetTag("msg-id")
	if found {
		switch msgID {
		case "subgift":
			{
				service.subHandler(client, message)
			}
		case "room_mods":
			{
				commaIndex := strings.Index(message.Params[1], ":")
				if commaIndex != -1 {
					mods := strings.Split(message.Params[1][commaIndex+2:], ", ")
					channel := strings.ToLower(message.Params[0][1:])
					service.modHandler(&channel, mods)
				}
			}
		case "resub":
			{
				service.subHandler(client, message)
			}

		case "sub":
			{
				service.subHandler(client, message)
			}
		}
	}

	if message.Command == "CLEARCHAT" {
		channelID := message.Tags["room-id"].Encode()
		channelInfo, _ := service.channelInfoService.GetChannelInfo(&channelID)
		if pubsub.IsWorking == false || channelInfo.ExtendedBansLogging == false {
			banDuration, banDurationFound := message.Tags.GetTag("ban-duration")
			intBanDuration := 0
			if banDurationFound {
				parsedValue, parseError := strconv.Atoi(banDuration)
				if parseError == nil {
					intBanDuration = parsedValue
				}
			}
			banReason, _ := message.Tags.GetTag("ban-reason")
			user := message.Params[1]
			channel := message.Params[0]
			messageType := "timeout"
			if intBanDuration == 0 {
				messageType = "ban"
			}
			formedMessage := chatMessage.ChatMessage{
				MessageStruct: chatMessage.MessageStruct{
					Date:        time.Now(),
					MessageType: messageType,
					BanLength:   intBanDuration,
					BanReason:   banReason},
				Channel:   channel,
				ChannelID: message.Tags["room-id"].Encode(),
				User:      user,
				UserID:    message.Tags["target-user-id"].Encode()}
			service.channelLogsService.Log(&formedMessage)
		} else {
			logger.Debug("Not logging")

		}
	}
	if message.Command == "PRIVMSG" {
		//	logger.Debug("Got PRIVMSG, parsing")

		formedMessage := chatMessage.ChatMessage{
			MessageStruct: chatMessage.MessageStruct{
				Username:    message.User,
				MessageType: "message",
				MessageBody: strings.TrimSpace(message.Params[1]),
				Date:        time.Now()},
			Channel:   message.Params[0][1:],
			ChannelID: message.Tags["room-id"].Encode(),
			User:      message.User,
			UserID:    message.Tags["user-id"].Encode(),
			IsMod:     message.Tags["mod"] == "1" || message.User == "khadesru" || message.Params[0][1:] == message.User,
			IsSub:     message.Tags["subscriber"] == "1",
			IsPrime:   strings.Contains(message.Tags["badges"].Encode(), "premium/1")}
		service.channelLogsService.Log(&formedMessage)
		channelInfo, _ := service.channelInfoService.GetChannelInfo(&formedMessage.ChannelID)
		service.autoMessageService.Decrement(channelInfo)

		commandBody, isCommand := formedMessage.GetCommand()

		isVote := strings.HasPrefix(message.Params[1], "%")
		if isVote == true {
			service.vote(client, channelInfo, &formedMessage)
		}
		// bits, bitsFound := message.Tags.GetTag("bits")
		// if bitsFound {
		// 	parsedBits, parsedBitsError := strconv.Atoi(bits)
		// 	if parsedBitsError == nil {
		// 		repos.AddBitsToUser(&formedMessage.ChannelID, &formedMessage.UserID, &formedMessage.User, parsedBits, "bits")
		// 	}
		// }
		if isCommand == true {
			logger.Debug("PRIVMGS is chat command")
			if message.User == "khadesru" && commandBody.Command == "debugsub" {
				subPlan := "2000"
				service.sendSubMessage(client, channelInfo, &formedMessage.User, &subPlan)
			}

			if message.User == "khadesru" && commandBody.Command == "debugresub" {
				resubCount := 3
				subPlan := "2000"
				service.sendResubMessage(client, channelInfo, &formedMessage.User, &resubCount, &subPlan)
				loggedSubscription := subscriptionInfo.SubscriptionInfo{
					User:      "khades",
					UserID:    "khades",
					ChannelID: formedMessage.ChannelID,
					Count:     resubCount,
					IsPrime:   false,
					SubPlan:   "1000",
					Date:      time.Now()}

				service.subscriptionInfoService.Log(&loggedSubscription)

				// service.eventBus.Publish(eventbus.EventSub(&formedMessage.ChannelID), "newsub")
				// service.eventBus.Publish(eventbus.Subtrain(&formedMessage.ChannelID), "newsub")
			}

			//handlerFunction := commandhandlers.Router.Go(commandBody.Command)
			//logger.Debug("Getting channel info")
			//handlerFunction(channelInfo, &formedMessage, commandBody, IrcClientInstance)
		}
	}
}
