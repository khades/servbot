package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/belak/irc"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

var chatHandler irc.HandlerFunc = func(client *irc.Client, message *irc.Message) {
	//fmt.Println(message.String())
	if message.Tags["msg-id"] == "room_mods" {
		commaIndex := strings.Index(message.Params[1], ":")
		if commaIndex != -1 {
			mods := strings.Split(message.Params[1][commaIndex+1:], ", ")
			for _, value := range mods {
				fmt.Println(value)

			}
		}
	}

	if message.Command == "PRIVMSG" {
		formedMessage := models.ChatMessage{
			Channel:     message.Params[0],
			Username:    message.User,
			MessageBody: message.Params[1],
			IsMod:       message.Tags["mod"] == "1",
			IsSub:       message.Tags["subscriber"] == "1",
			Date:        time.Now()}
		repos.LogMessage(formedMessage)
	}
	if message.Command == "001" {
		client.Write("CAP REQ twitch.tv/tags")
		client.Write("CAP REQ twitch.tv/membership")
		client.Write("CAP REQ twitch.tv/commands")
		for _, value := range repos.Config.Channels {
			client.Write("JOIN #" + value)
		}
		client.Write("PRIVMSG #nuke73 .mods")
	}
}
