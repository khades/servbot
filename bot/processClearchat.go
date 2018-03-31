package bot

import (
	"strconv"
	"time"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	irc "gopkg.in/irc.v2"
)

func processClearchat(message *irc.Message) {

	banDuration, banDurationFound := message.Tags.GetTag("ban-duration")
	intBanDuration := 0
	if banDurationFound {
		parsedValue, parseError := strconv.Atoi(banDuration)
		if parseError == nil {
			intBanDuration = parsedValue
		}
	}
	banReason, _ := message.Tags.GetTag("ban-reason")
	if len(message.Params) < 2 {
		return
	}
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
