package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func resubHandler(message *irc.Message, ircClient *ircClient.IrcClient) {
	msgParamMonths, msgParamMonthsFound := message.Tags.GetTag("msg-param-months")
	user, userFound := message.Tags.GetTag("display-name")
	channel := message.Params[0][1:]
	if msgParamMonthsFound && userFound && channel != "" {
		resubCount, resubCountError := strconv.Atoi(msgParamMonths)
		if resubCountError == nil {
			formedMessage := models.ChatMessage{
				MessageStruct: models.MessageStruct{
					Date:     time.Now(),
					SubCount: resubCount},
				Channel: channel,
				User:    user,
				IsPrime: strings.Contains(message.String(), "Twitch Prime"),
			}
			repos.LogMessage(&formedMessage)
			sendResubMessage(&channel, &user, &formedMessage.SubCount)
			log.Printf("Channel %v: %v resubbed for %v months\n", formedMessage.Channel, formedMessage.User, formedMessage.SubCount)
		}
	}
}
