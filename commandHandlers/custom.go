package commandHandlers

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode/utf8"

	"html"

	"github.com/hoisie/mustache"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func Short(s string, i int) string {
	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}

// Custom handler checks if input command has template and then fills it in with mustache templating and sends to a specified/user
func Custom(online bool, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {
	channelInfo, error := repos.GetChannelInfo(&chatMessage.ChannelID)
	if error != nil {
		return
	}
	channelStatus := &models.ChannelInfoForTemplate{ChannelInfo: *channelInfo, IsMod: chatMessage.IsMod}
	channelStatus.IsMod = chatMessage.IsMod
	template, err := repos.GetChannelTemplate(&chatMessage.ChannelID, &chatCommand.Command)
	user := chatMessage.User
	if err != nil || template.Template == "" {
		return
	}

	if chatMessage.IsMod == false {
		if channelStatus.StreamStatus.Online == true && template.ShowOnline == false {
			return
		}
		if channelStatus.StreamStatus.Online == false && template.ShowOffline == false {
			return
		}
	}
	if template.IntegerRandomizer.Enabled == true && template.IntegerRandomizer.UpperLimit > template.IntegerRandomizer.LowerLimit {
		channelStatus.RandomInteger = template.IntegerRandomizer.LowerLimit + rand.Intn(template.IntegerRandomizer.UpperLimit-template.IntegerRandomizer.LowerLimit)
		if channelStatus.RandomInteger == template.IntegerRandomizer.LowerLimit {
			channelStatus.RandomIntegerIsMinimal = true
		}
		if channelStatus.RandomInteger == template.IntegerRandomizer.UpperLimit {
			channelStatus.RandomIntegerIsMaximal = true
		}
		if channelStatus.RandomInteger == 0 {
			channelStatus.RandomIntegerIsMinimal = true
		}
		if template.IntegerRandomizer.TimeoutAfter == true && channelStatus.RandomInteger > 0 {
			if chatMessage.IsMod == false {
				ircClient.SendPublic(&models.OutgoingMessage{
					Channel: chatMessage.Channel,
					Body:    fmt.Sprintf("/timeout %s %d ", user, channelStatus.RandomInteger)})
			} else {
				ircClient.SendPublic(&models.OutgoingMessage{
					Channel: chatMessage.Channel,
					User:    user,
					Body:    "Модератора нельзя затаймаутить SMOrc"})
				return
			}

		}
	}

	if template.StringRandomizer.Enabled == true {
		if len(template.StringRandomizer.Strings) == 0 {
			commandValues := strings.Split(chatCommand.Body, ",")
			if len(commandValues) != 0 {
				channelStatus.RandomString = strings.TrimSpace(commandValues[rand.Intn(len(commandValues)-1)])

			}
		} else {
			channelStatus.RandomString = strings.TrimSpace(template.StringRandomizer.Strings[rand.Intn(len(template.StringRandomizer.Strings)-1)])
		}
	}

	message := mustache.Render(template.Template, channelStatus)
	if utf8.RuneCountInString(message) > 400 {
		message = Short(message, 397) + "..."
	}
	redirectTo := chatMessage.User
	if chatCommand.Body != "" && !(template.StringRandomizer.Enabled == true && len(template.StringRandomizer.Strings) == 0) && template.PreventRedirect == false {
		if strings.HasPrefix(chatCommand.Body, "@") {
			redirectTo = chatCommand.Body[1:]
		} else {
			redirectTo = chatCommand.Body

		}
	}
	if template.PreventDebounce == true {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    user,
			Body:    html.UnescapeString(message)})
	} else {
		ircClient.SendDebounced(models.OutgoingDebouncedMessage{
			Message: models.OutgoingMessage{
				Channel: chatMessage.Channel,
				User:    user,
				Body:    html.UnescapeString(message)},
			Command:    template.AliasTo,
			RedirectTo: redirectTo})
	}

}
