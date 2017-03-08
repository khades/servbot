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
	msgID, found := message.Tags.GetTag("msg-id")
	if found {
		switch msgID {
		case "room_mods":
			{
				commaIndex := strings.Index(message.Params[1], ":")
				log.Println(message.String())
				if commaIndex != -1 {
					mods := strings.Split(message.Params[1][commaIndex+2:], ", ")
					channel := message.Params[0][1:]
					modHandler(&channel, &mods)
				}
			}
		case "resub":
			{
				resubHandler(message, &IrcClientInstance)
			}
		}
	}

	if message.User == "twitchnotify" {
		subHandler(message, &IrcClientInstance)
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
		formedMessage := models.ChatMessage{
			MessageStruct: models.MessageStruct{
				Date:      time.Now(),
				BanLength: intBanDuration,
				BanReason: banReason},
			Channel:   channel,
			ChannelID: message.Tags["room-id"].Encode(),
			User:      user,
			UserID:    message.Tags["user-id"].Encode()}
		repos.LogMessage(&formedMessage)
		//	log.Printf("Channel %v: %v is banned for %v \n", channel, user, intBanDuration)
	}
	if message.Command == "PRIVMSG" {
		log.Println(message.String())
		log.Println(message.Tags["room-id"].Encode())
		formedMessage := models.ChatMessage{
			MessageStruct: models.MessageStruct{
				Username:    message.User,
				MessageBody: message.Params[1],
				Date:        time.Now()},
			Channel:   message.Params[0][1:],
			ChannelID: message.Tags["room-id"].Encode(),
			User:      message.User,
			UserID:    message.Tags["user-id"].Encode(),

			IsMod:   message.Tags["mod"] == "1" || message.User == "khadesru" || message.Params[0][1:] == message.User,
			IsSub:   message.Tags["subscriber"] == "1",
			IsPrime: strings.Contains(message.Tags["badges"].Encode(), "premium/1")}
		repos.LogMessage(&formedMessage)
		repos.DecrementAutoMessages(&formedMessage.ChannelID)
		commandBody, isCommand := formedMessage.GetCommand()
		if isCommand {
			if message.User == "khadesru" && commandBody.Command == "debugSub" {
				sendSubMessage(&formedMessage.Channel, &formedMessage.ChannelID, &formedMessage.User)
				log.Println("debug subalert is triggered")
			}
			if message.User == "khadesru" && commandBody.Command == "debugResub" {
				resubCount := 3
				sendResubMessage(&formedMessage.Channel, &formedMessage.ChannelID, &formedMessage.User, &resubCount)
				log.Println("debug resubalert is triggered")
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
		IrcClientInstance = ircClient.IrcClient{Client: client, Bounces: make(map[string]time.Time), Ready: true}
		IrcClientInstance.SendModsCommand()
		log.Println("Bot is started")
	}
}
