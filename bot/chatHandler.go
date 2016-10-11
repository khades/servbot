package bot

import (
	"fmt"
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
	//	log.Println(message.String())
	msgID, found := message.Tags.GetTag("msg-id")
	if found {
		switch msgID {
		case "room_mods":
			{
				commaIndex := strings.Index(message.Params[1], ":")
				if commaIndex != -1 {
					//				log.Printf("Channel %v: got mods list", message.Params[0])
					mods := strings.Split(message.Params[1][commaIndex+2:], ", ")
					repos.PushMods(message.Params[0][1:], mods)
				}
			}
		case "resub":
			{
				msgParamMonths, msgParamMonthsFound := message.Tags.GetTag("msg-param-months")
				user, userFound := message.Tags.GetTag("display-name")
				channel := message.Params[0][1:]
				if msgParamMonthsFound && userFound && channel != "" {
					resubCount, resubCountError := strconv.Atoi(msgParamMonths)
					if resubCountError == nil {
						formedMessage := models.ChatMessage{
							Channel:          channel,
							User:             user,
							IsMod:            false,
							IsSub:            true,
							Date:             time.Now(),
							SubscriptionInfo: &models.SubscriptionInfo{Count: resubCount}}
						repos.LogMessage(formedMessage)
						channelInfo, error := repos.GetChannelInfo(channel)
						if error == nil && channelInfo.SubAlert.Enabled == true {
							messageBody := strings.TrimSpace(fmt.Sprintf("%s %s%s",
								channelInfo.SubAlert.RepeatPrefix,
								strings.Repeat(channelInfo.SubAlert.RepeatBody+" ", formedMessage.SubscriptionInfo.Count),
								channelInfo.SubAlert.RepeatPostfix))
							if messageBody != "" {
								IrcClientInstance.SendPublic(models.OutgoingMessage{
									Body:    messageBody,
									Channel: channel,
									User:    user})
							}
						}
						log.Printf("Channel %v: %v resubbed for %v months\n", formedMessage.Channel, formedMessage.User, formedMessage.SubscriptionInfo.Count)
					}
				}
			}
		}
	}

	if message.User == "twitchnotify" {
		log.Println("Got first sub")
		user := strings.Split(message.Params[1], " ")[0]
		channel := message.Params[0][1:]
		if user != "" && channel != "" {
			formedMessage := models.ChatMessage{
				Channel:          channel,
				User:             user,
				IsMod:            false,
				IsSub:            true,
				IsPrime:          strings.Contains(message.String(), "Twitch Prime"),
				Date:             time.Now(),
				SubscriptionInfo: &models.SubscriptionInfo{Count: 1}}
			repos.LogMessage(formedMessage)
			channelInfo, error := repos.GetChannelInfo(channel)
			if error == nil && channelInfo.SubAlert.Enabled == true && channelInfo.SubAlert.FirstMessage != "" {
				IrcClientInstance.SendPublic(models.OutgoingMessage{
					Body:    channelInfo.SubAlert.FirstMessage,
					Channel: channel,
					User:    user})
			}
			log.Printf("Channel %v: %v subbed\n", formedMessage.Channel, formedMessage.User)
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
		formedMessage := models.ChatMessage{
			Channel: channel,
			User:    user,
			IsMod:   false,
			IsSub:   true,
			Date:    time.Now(),
			BanInfo: &models.BanInfo{Duration: intBanDuration, Reason: banReason}}
		repos.LogMessage(formedMessage)
		//	log.Printf("Channel %v: %v is banned for %v \n", channel, user, intBanDuration)
	}
	if message.Command == "PRIVMSG" {
		formedMessage := models.ChatMessage{
			Channel:     message.Params[0][1:],
			User:        message.User,
			MessageBody: message.Params[1],
			IsMod:       message.Tags["mod"] == "1" || message.User == "khadesru",
			IsSub:       message.Tags["subscriber"] == "1",
			IsPrime:     strings.Contains(message.Tags["badges"].Encode(), "premium/1"),
			Date:        time.Now()}
		repos.LogMessage(formedMessage)
		isCommand, commandBody := formedMessage.IsCommand()
		if isCommand {
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
