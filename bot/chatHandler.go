package bot

import (
	"strconv"
	"strings"
	"time"

	"github.com/khades/servbot/commandhandlers"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/pubsub"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

var chatHandler irc.HandlerFunc = func(client *irc.Client, message *irc.Message) {
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
				subHandler(message, IrcClientInstance)
			}
		case "room_mods":
			{
				commaIndex := strings.Index(message.Params[1], ":")
				if commaIndex != -1 {
					mods := strings.Split(message.Params[1][commaIndex+2:], ", ")
					channel := strings.ToLower(message.Params[0][1:])
					modHandler(&channel, mods)
				}
			}
		case "resub":
			{
				subHandler(message, IrcClientInstance)
			}

		case "sub":
			{
				subHandler(message, IrcClientInstance)
			}
		}
	}

	if message.Command == "CLEARCHAT" {
		channelID := message.Tags["room-id"].Encode()
		channelInfo, _ := repos.GetChannelInfo(&channelID)
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
			formedMessage := models.ChatMessage{
				MessageStruct: models.MessageStruct{
					Date:        time.Now(),
					MessageType: messageType,
					BanLength:   intBanDuration,
					BanReason:   banReason},
				Channel:   channel,
				ChannelID: message.Tags["room-id"].Encode(),
				User:      user,
				UserID:    message.Tags["target-user-id"].Encode()}
			repos.LogMessage(&formedMessage)
		} else {
			logger.Debug("Not logging")

		}
	}
	if message.Command == "PRIVMSG" {
		//	logger.Debug("Got PRIVMSG, parsing")

		formedMessage := models.ChatMessage{
			MessageStruct: models.MessageStruct{
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
		//	logger.Debug("Logging PRIVMSG")
		repos.LogMessage(&formedMessage)
		channelInfo, _ := repos.GetChannelInfo(&formedMessage.ChannelID)
		repos.DecrementAutoMessages(channelInfo)

		commandBody, isCommand := formedMessage.GetCommand()

		isVote := strings.HasPrefix(message.Params[1], "%")
		if isVote == true {

			commandhandlers.Vote(channelInfo, &formedMessage, commandBody, IrcClientInstance)
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
			logger.Debug(message.User)
			logger.Debug(commandBody.Command)
			if message.User == "khadesru" && commandBody.Command == "debugsub" {
				subPlan := "2000"
				sendSubMessage(channelInfo, &formedMessage.User, &subPlan)

			}
			if message.User == "khadesru" && commandBody.Command == "debugresub" {
				resubCount := 3
				subPlan := "2000"
				sendResubMessage(channelInfo, &formedMessage.User, &resubCount, &subPlan)
				loggedSubscription := models.SubscriptionInfo{
					User:      "khades",
					UserID:    "khades",
					ChannelID: formedMessage.ChannelID,
					Count:     resubCount,
					IsPrime:   false,
					SubPlan:   "1000",
					Date:      time.Now()}

				repos.LogSubscription(&loggedSubscription)

				eventbus.EventBus.Publish(eventbus.EventSub(&formedMessage.ChannelID), "newsub")
				eventbus.EventBus.Publish(eventbus.Subtrain(&formedMessage.ChannelID), "newsub")
			}
			handlerFunction := commandhandlers.Router.Go(commandBody.Command)
			logger.Debug("Getting channel info")
			handlerFunction(channelInfo, &formedMessage, commandBody, IrcClientInstance)
		}
	}

	if message.Command == "001" {
		client.Write("CAP REQ twitch.tv/tags")
		client.Write("CAP REQ twitch.tv/membership")
		client.Write("CAP REQ twitch.tv/commands")
		activeChannels, _ := repos.GetActiveChannels()
		for _, value := range activeChannels {
			client.Write("JOIN #" + value.Channel)
		}
		IrcClientInstance = &ircClient.IrcClient{Client: client, Bounces: make(map[string]time.Time), Ready: true, ModChannelIndex: 0, MessageQueue: []string{}}
		IrcClientInstance.SendModsCommand()
		logger.Info("Bot is started")
	}
}
