package commandhandlers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"html"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func short(s string, i int) string {
	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}

type templateExtendedObject struct {
	models.ChannelInfo
	UserID                 string
	RandomInteger          int
	RandomIntegerIsMinimal bool
	RandomIntegerIsMaximal bool
	RandomIngegerIsZero    bool
	RandomString           string
	IsMod                  bool
	IsSub                  bool
	CommandBody            string
	CommandBodyIsEmpty     bool
}
type FollowerDuration struct {
	IsFollower       bool
	FollowerDuration string
}

func (channelInfo templateExtendedObject) FollowerInfo() FollowerDuration {

	isFollower, dur := repos.GetIfFollowerToChannel(&channelInfo.ChannelID, &channelInfo.UserID)
	if isFollower == true {
		return FollowerDuration{true, l10n.HumanizeDuration(time.Now().Sub(dur), channelInfo.Lang)}

	}
	return FollowerDuration{false, ""}
}

// custom handler checks if input command has template and then fills it in with mustache templating and sends to a specified/user
func custom(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {

	templateObject := &templateExtendedObject{ChannelInfo: *channelInfo, IsMod: chatMessage.IsMod, CommandBody: chatCommand.Body, CommandBodyIsEmpty: chatCommand.Body == ""}
	templateObject.IsMod = chatMessage.IsMod
	templateObject.IsSub = chatMessage.IsSub
	templateObject.UserID = chatMessage.UserID
	template, err := repos.GetChannelTemplate(&chatMessage.ChannelID, &chatCommand.Command)
	user := chatMessage.User
	if err != nil || template.Template == "" {
		return
	}

	if chatMessage.IsMod == false {
		if templateObject.StreamStatus.Online == true && template.ShowOnline == false {
			return
		}
		if templateObject.StreamStatus.Online == false && template.ShowOffline == false {
			return
		}
	}
	if template.IntegerRandomizer.Enabled == true && template.IntegerRandomizer.UpperLimit > template.IntegerRandomizer.LowerLimit {
		templateObject.RandomInteger = template.IntegerRandomizer.LowerLimit + rand.Intn(template.IntegerRandomizer.UpperLimit-template.IntegerRandomizer.LowerLimit)
		if templateObject.RandomInteger == template.IntegerRandomizer.LowerLimit {
			templateObject.RandomIntegerIsMinimal = true
		}
		if templateObject.RandomInteger == template.IntegerRandomizer.UpperLimit {
			templateObject.RandomIntegerIsMaximal = true
		}
		if templateObject.RandomInteger == 0 {
			templateObject.RandomIntegerIsMinimal = true
		}
		if template.IntegerRandomizer.TimeoutAfter == true && templateObject.RandomInteger > 0 {
			if chatMessage.IsMod == false {
				ircClient.SendPublic(&models.OutgoingMessage{
					Channel: chatMessage.Channel,
					Body:    fmt.Sprintf("/timeout %s %d ", user, templateObject.RandomInteger)})
			} else {
				ircClient.SendPublic(&models.OutgoingMessage{
					Channel: chatMessage.Channel,
					User:    user,
					Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CantMuteModerator})
				return
			}

		}
	}

	if template.StringRandomizer.Enabled == true {
		if len(template.StringRandomizer.Strings) == 0 {
			commandValues := strings.Split(chatCommand.Body, ",")
			if len(commandValues) != 0 {
				templateObject.RandomString = strings.TrimSpace(commandValues[rand.Intn(len(commandValues)-1)])

			}
		} else {
			templateObject.RandomString = strings.TrimSpace(template.StringRandomizer.Strings[rand.Intn(len(template.StringRandomizer.Strings)-1)])
		}
	}

	message, renderError := mustache.Render(template.Template, templateObject)
	if renderError != nil {
		return
	}
	message = strings.TrimSpace(message)
	if utf8.RuneCountInString(message) > 400 {
		message = short(message, 397) + "..."
	}
	if message == "" {
		return
	}
	redirectTo := chatMessage.User
	if chatCommand.Body != "" && !(template.StringRandomizer.Enabled == true && len(template.StringRandomizer.Strings) == 0) && template.PreventRedirect == false {
		if strings.HasPrefix(chatCommand.Body, "@") {
			redirectTo = chatCommand.Body[1:]
		} else {
			redirectTo = chatCommand.Body

		}
	}
	redirectTo = strings.Replace(redirectTo, "@", " @", -1)
	if template.OnlyPrivate == true {
		ircClient.SendPrivate(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    user,
			Body:    html.UnescapeString(message)})
		return
	}

	if template.PreventDebounce == true {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    user,
			Body:    html.UnescapeString(message)})
		return
	}

	ircClient.SendDebounced(models.OutgoingDebouncedMessage{
		Message: models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    user,
			Body:    html.UnescapeString(message)},
		Command:    template.AliasTo,
		RedirectTo: redirectTo})

}
