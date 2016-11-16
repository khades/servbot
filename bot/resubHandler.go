package bot

import (
	"log"
	"strconv"
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
				Channel:          channel,
				User:             user,
				IsMod:            false,
				IsSub:            true,
				Date:             time.Now(),
				SubscriptionInfo: &models.SubscriptionInfo{Count: resubCount}}
			repos.LogMessage(&formedMessage)
			sendResubMessage(&channel, &user, &formedMessage.SubscriptionInfo.Count)
			log.Printf("Channel %v: %v resubbed for %v months\n", formedMessage.Channel, formedMessage.User, formedMessage.SubscriptionInfo.Count)
		}
	}
}
