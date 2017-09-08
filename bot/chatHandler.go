package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/commandHandlers"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

var chatHandler irc.HandlerFunc = func(client *irc.Client, message *irc.Message) {
	if strings.Contains(message.String(), ":jtv") {
		log.Println("JTV: " + message.String())
	}
	if strings.Contains(message.String(), ":tmi.twitch.tv") {
		log.Println("TMI.TWITCH.TV: " + message.String())
	}
	msgID, found := message.Tags.GetTag("msg-id")
	if found {
		switch msgID {
		case "room_mods":
			{
				commaIndex := strings.Index(message.Params[1], ":")
				if commaIndex != -1 {
					mods := strings.Split(message.Params[1][commaIndex+2:], ", ")
					channel := message.Params[0][1:]
					modHandler(&channel, &mods)
				}
			}
		case "resub":
			{
				subHandler(message, &IrcClientInstance)
			}

		case "sub":
			{
				subHandler(message, &IrcClientInstance)
			}
		}
	}

	if message.Command == "CLEARCHAT" {
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
	}
	if message.Command == "PRIVMSG" {
		formedMessage := models.ChatMessage{
			MessageStruct: models.MessageStruct{
				Username:    message.User,
				MessageType: "message",
				MessageBody: message.Params[1],
				Date:        time.Now()},
			Channel:   message.Params[0][1:],
			ChannelID: message.Tags["room-id"].Encode(),
			User:      message.User,
			UserID:    message.Tags["user-id"].Encode(),
			IsMod:     message.Tags["mod"] == "1" || message.User == "khadesru" || message.Params[0][1:] == message.User,
			IsSub:     message.Tags["subscriber"] == "1",
			IsPrime:   strings.Contains(message.Tags["badges"].Encode(), "premium/1")}
		repos.LogMessage(&formedMessage)
		repos.DecrementAutoMessages(&formedMessage.ChannelID)
		commandBody, isCommand := formedMessage.GetCommand()
		bits, bitsFound := message.Tags.GetTag("bits")
		if bitsFound {
			parsedBits, parsedBitsError := strconv.Atoi(bits)
			if parsedBitsError == nil {
				repos.AddBitsToUser(&formedMessage.ChannelID, &formedMessage.UserID, &formedMessage.User, parsedBits, "bits")
			}
		}
		if isCommand {
			if message.User == "khadesru" && commandBody.Command == "debugSub" {
				subPlan := "2000"
				sendSubMessage(&formedMessage.Channel, &formedMessage.ChannelID, &formedMessage.User, &subPlan)
			}
			if message.User == "khadesru" && commandBody.Command == "debugResub" {
				resubCount := 3
				subPlan := "2000"
				sendResubMessage(&formedMessage.Channel, &formedMessage.ChannelID, &formedMessage.User, &resubCount, &subPlan)
			}
			handlerFunction := commandHandlers.Router.Go(commandBody.Command)
			handlerFunction(true, &formedMessage, commandBody, &IrcClientInstance)
		}
	}

	if message.Command == "001" {
		client.Write("CAP REQ twitch.tv/tags")
		client.Write("CAP REQ twitch.tv/membership")
		client.Write("CAP REQ twitch.tv/commands")
		for _, value := range repos.Config.Channels {
			client.Write("JOIN #" + value)
		}
		IrcClientInstance = ircClient.IrcClient{Client: client, Bounces: make(map[string]time.Time), Ready: true, ModChannelIndex: 0}
		IrcClientInstance.SendModsCommand()
		log.Println("Bot is started")
	}
}
