package twitchIRCHandler

import (
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/followers"
	"github.com/khades/servbot/songRequest"
	"github.com/khades/servbot/subscriptionInfo"
	"github.com/khades/servbot/template"

	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/subday"
	"github.com/khades/servbot/twitchIRC"
	"github.com/khades/servbot/userResolve"

	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/pubsub"
	"gopkg.in/irc.v2"
)

type TwitchIRCHandler struct {
	subdayService           *subday.Service
	channelInfoService      *channelInfo.Service
	subAlertService         *subAlert.Service
	channelLogsService      *channelLogs.Service
	autoMessageService      *autoMessage.Service
	userResolveService      *userResolve.Service
	subscriptionInfoService *subscriptionInfo.Service
	templateService         *template.Service
	followersService        *followers.Service
	songRequestService      *songRequest.Service
	eventBus                EventBus.Bus
	pubsub                  *pubsub.Client
	// eventService            *event.Service
	// balanceService          *balance.Service
}

func (service *TwitchIRCHandler) Handle(client *twitchIRC.Client, message *irc.Message) {
	//logger := logrus.WithFields(logrus.Fields{
	//	"package": "twitchIRCHandler",
	//	"action":  "Handle"})
	//if strings.Contains(message.String(), ":jtv") {
	//	logger.Debug("JTV: " + message.String())
	//}
	//if strings.Contains(message.String(), ":tmi.twitch.tv") {
	//	logger.Debug("TMI.TWITCH.TV: " + message.String())
	//}
	msgID, found := message.Tags.GetTag("msg-id")
	if found {
		switch msgID {
		case "subgift":
			{
				service.sub(client, message)
			}
		case "anonsubgift":
			{
				service.sub(client, message)
			}
		case "room_mods":
			{
				commaIndex := strings.Index(message.Params[1], ":")
				if commaIndex != -1 {
					mods := strings.Split(message.Params[1][commaIndex+2:], ", ")
					channel := strings.ToLower(message.Params[0][1:])
					service.mod(&channel, mods)
				}
			}
		case "resub":
			{
				service.sub(client, message)
			}

		case "sub":
			{
				service.sub(client, message)
			}
		}
	}

	if message.Command == "CLEARCHAT" {
		channelID := message.Tags["room-id"].Encode()
		channelInfo, _ := service.channelInfoService.Get(&channelID)
		if service.pubsub.IsWorking == false || channelInfo.ExtendedBansLogging == false {
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
		}

		//else {
		//	logger.Debug("Not logging")
		//
		//}
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
		channelInfo, _ := service.channelInfoService.Get(&formedMessage.ChannelID)
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

		// 		service.balanceService.Inc(
		// 			formedMessage.ChannelID,
		// 			formedMessage.UserID,
		// 			formedMessage.User,
		// 			float64(parsedBits))
		// 		service.eventService.Put(formedMessage.ChannelID, event.Event{
		// 			User:     message.User,
		// 			Type:     event.BITS,
		// 			Amount:   parsedBits,
		// 			Message:  formedMessage.MessageStruct.MessageBody,
		// 			Currency: "USD",
		// 		})
		// 	}
		// }
		if isCommand == true {
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

				service.eventBus.Publish(eventbus.EventSub(&formedMessage.ChannelID), "newsub")
				service.eventBus.Publish(eventbus.Subtrain(&formedMessage.ChannelID), "newsub")
			}
			if commandBody.Command == "new" {
				service.new(channelInfo, &formedMessage, commandBody, client)
			} else if commandBody.Command == "alias" {
				service.alias(channelInfo, &formedMessage, commandBody, client)
			} else if commandBody.Command == "subdaynew" {
				service.subdayNew(channelInfo, &formedMessage, commandBody, client)
			} else {
				service.custom(channelInfo, &formedMessage, commandBody, client)
			}

		}
	}
}
