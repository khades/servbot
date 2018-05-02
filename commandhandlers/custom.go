package commandhandlers

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
	UserID string
	User   string
	// RandomInteger          int
	// RandomIntegerIsMinimal bool
	// RandomIntegerIsMaximal bool
	// RandomIngegerIsZero    bool
	// RandomString           string
	IsMod              bool
	IsSub              bool
	CommandBody        string
	CommandBodyIsEmpty bool
	PreventDebounce    bool
	PreventRedirect    bool
	IrcClient          *ircClient.IrcClient
}
type songPullResult struct {
	Success     bool
	PulledVideo models.SongRequest
}
type followerDuration struct {
	IsFollower       bool
	FollowerDuration string
}
type banmeResult struct {
	Banned      bool
	Moderator   bool
	BanDuration int
}

func (channelInfo *templateExtendedObject) FollowerInfo() followerDuration {
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	isFollower, dur := repos.GetIfFollowerToChannel(&channelInfo.ChannelID, &channelInfo.UserID)
	if isFollower == true {
		return followerDuration{true, l10n.HumanizeDuration(time.Now().Sub(dur), channelInfo.Lang)}

	}
	return followerDuration{false, ""}
}

func (channelInfo templateExtendedObject) CurrentSong() models.CurrentSong {
	return repos.GetTopRequest(&channelInfo.ChannelID, channelInfo.Lang)
}

func (channelInfo *templateExtendedObject) Random() string {
	lowerLimit := 0
	upperLimit := 100
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	if channelInfo.IsMod == false && channelInfo.StreamStatus.Online == true {
		return ""
	}
	return strconv.Itoa(lowerLimit + rand.Intn(upperLimit-lowerLimit))
}

func (channelInfo *templateExtendedObject) Banme30() banmeResult {
	return channelInfo.Banme(30)
}

func (channelInfo *templateExtendedObject) Banme60() banmeResult {
	return channelInfo.Banme(60)
}

func (channelInfo *templateExtendedObject) Banme300() banmeResult {
	return channelInfo.Banme(300)
}

func (channelInfo *templateExtendedObject) Banme600() banmeResult {
	return channelInfo.Banme(600)
}

func (channelInfo *templateExtendedObject) SubdayEnd() string {
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	if channelInfo.IsMod == false {
		return ""
	}
	if channelInfo.SubdayIsActive == true {
		repos.CloseActiveSubday(&channelInfo.ChannelID, &channelInfo.User, &channelInfo.UserID)
		return l10n.GetL10n(channelInfo.GetChannelLang()).SubdayEndSuccess
	}

	return l10n.GetL10n(channelInfo.GetChannelLang()).SubdayEndNothingToClose
}

func (channelInfo *templateExtendedObject) Banme(length int) banmeResult {
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	if channelInfo.IsMod == true {
		return banmeResult{
			Banned:    false,
			Moderator: true}
	}
	banDuration := rand.Intn(length)
	if banDuration == 0 {
		return banmeResult{
			Banned: false}
	}

	channelInfo.IrcClient.SendPublic(&models.OutgoingMessage{
		Channel: channelInfo.Channel,
		Body:    fmt.Sprintf("/timeout %s %d ", channelInfo.User, banDuration)})

	return banmeResult{
		Banned:      true,
		BanDuration: banDuration}
}

func (channelInfo *templateExtendedObject) Pick() string {
	if channelInfo.IsMod == false && channelInfo.StreamStatus.Online == true {
		return ""
	}
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	commandValues := strings.Split(channelInfo.CommandBody, ",")
	if len(commandValues) == 1 {
		return strings.TrimSpace(commandValues[0])
	}
	if len(commandValues) != 0 {
		return strings.TrimSpace(commandValues[rand.Intn(len(commandValues))])

	}

	return "SMOrc"
}
func (channelInfo *templateExtendedObject) Ask() string {
	if channelInfo.IsMod == false && channelInfo.StreamStatus.Online == true {
		return ""
	}
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	variants := []string{
		"Бесспорно",
		"Предрешено",
		"Никаких сомнений MingLee",
		"Определённо да VoHiYo",
		"Можешь быть уверен в этом Keepo",
		"Мне кажется — «да»",
		"Вероятнее всего",
		"Хорошие перспективы Keepo",
		"Знаки говорят — «да»",
		"Да",
		"Пока не ясно, попробуй снова",
		"Спроси позже ResidentSleeper",
		"Лучше не рассказывать 4Head",
		"Сейчас нельзя предсказать ResidentSleeper",
		"Сконцентрируйся и спроси опять",
		"Даже не думай WutFace",
		"Мой ответ — «нет» SMOrc",
		"По моим данным — «нет»",
		"Перспективы не очень хорошие",
		"Весьма сомнительно SMOrc"}
	return strings.TrimSpace(variants[rand.Intn(len(variants))])
}

func (channelInfo *templateExtendedObject) AddSongRequest() string {
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	result := repos.AddSongRequest(&channelInfo.User, channelInfo.IsSub, &channelInfo.UserID, &channelInfo.ChannelID, &channelInfo.CommandBody)
	if result.Success == true {
		result.LengthStr = l10n.HumanizeDuration(result.Length, channelInfo.Lang)
	}

	if result.YoutubeRestricted == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestYoutubeRestricted, result.Title)
	}

	if result.TwitchRestricted == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestTwitchRestricted, result.Title)
	}

	if result.ChannelRestricted == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestChannelRestricted, result.Title)
	}

	if result.TagRestricted == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestTagRestricted, result.Title, result.Tag)
	}

	if result.Offline == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestOffline)
	}

	if result.PlaylistIsFull == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPlaylistIsFull)
	}

	if result.AlreadyInPlaylist == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestAlreadyInPlaylist, result.Title)
	}

	if result.TooManyRequests == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestTooManyRequests)
	}

	if result.InvalidLink == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestInvalidLink)
	}

	if result.NothingFound == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestNothingFound)
	}

	if result.InternalError == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestInternalError)
	}

	if result.TooLong == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestTooLong, result.Title)
	}

	if result.TooLittleViews == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestTooLittleViews, result.Title)
	}

	if result.MoreDislikes == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestMoreDislikes, result.Title)
	}

	if result.Success == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestSuccess, result.Title, result.LengthStr)
	}
	return l10n.GetL10n(channelInfo.Lang).SongRequestInternalError
}

func (channelInfo *templateExtendedObject) SetSongRequestVolume() string {
	channelInfo.PreventRedirect = true
	if channelInfo.IsMod == false {
		return l10n.GetL10n(channelInfo.Lang).SongRequestNotAModerator
	}
	channelInfo.PreventDebounce = true

	volume, volumeError := strconv.ParseInt(channelInfo.CommandBody, 10, 23)
	if volumeError != nil {
		return l10n.GetL10n(channelInfo.Lang).VolumeChangeInvalidValue
	}
	if volume > 100 || volume < 0 {
		return l10n.GetL10n(channelInfo.Lang).VolumeChangeInvalidValue
	}
	repos.SetSongRequestVolume(&channelInfo.ChannelID, int(volume))
	return l10n.GetL10n(channelInfo.Lang).VolumeChangeSuccess
}

func (channelInfo *templateExtendedObject) PullSongRequest() string {
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true

	pulledVideo, pulled := repos.PullLastUserSongRequest(&channelInfo.ChannelID, &channelInfo.UserID)
	if pulled == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulled, pulledVideo.Title)
	}
	return l10n.GetL10n(channelInfo.Lang).SongRequestNoRequests
}

// custom handler checks if input command has template and then fills it in with mustache templating and sends to a specified/user
func custom(channelInfo *models.ChannelInfo, chatMessage *models.ChatMessage, chatCommand models.ChatCommand, ircClient *ircClient.IrcClient) {

	templateObject := &templateExtendedObject{ChannelInfo: *channelInfo, IsMod: chatMessage.IsMod, CommandBody: chatCommand.Body, CommandBodyIsEmpty: chatCommand.Body == ""}
	templateObject.IsMod = chatMessage.IsMod
	templateObject.IsSub = chatMessage.IsSub
	templateObject.UserID = chatMessage.UserID
	templateObject.User = chatMessage.User
	templateObject.IrcClient = ircClient
	template, err := repos.GetChannelTemplate(&chatMessage.ChannelID, &chatCommand.Command)

	if err != nil || template.Template == "" {
		return
	}

	// if template.IntegerRandomizer.Enabled == true && template.IntegerRandomizer.UpperLimit > template.IntegerRandomizer.LowerLimit {
	// 	templateObject.RandomInteger = template.IntegerRandomizer.LowerLimit + rand.Intn(template.IntegerRandomizer.UpperLimit-template.IntegerRandomizer.LowerLimit)
	// 	if templateObject.RandomInteger == template.IntegerRandomizer.LowerLimit {
	// 		templateObject.RandomIntegerIsMinimal = true
	// 	}
	// 	if templateObject.RandomInteger == template.IntegerRandomizer.UpperLimit {
	// 		templateObject.RandomIntegerIsMaximal = true
	// 	}
	// 	if templateObject.RandomInteger == 0 {
	// 		templateObject.RandomIntegerIsMinimal = true
	// 	}
	// 	if template.IntegerRandomizer.TimeoutAfter == true && templateObject.RandomInteger > 0 {
	// 		if chatMessage.IsMod == false {
	// 			ircClient.SendPublic(&models.OutgoingMessage{
	// 				Channel: chatMessage.Channel,
	// 				Body:    fmt.Sprintf("/timeout %s %d ", templateObject.User, templateObject.RandomInteger)})
	// 		} else {
	// 			ircClient.SendPublic(&models.OutgoingMessage{
	// 				Channel: chatMessage.Channel,
	// 				User:    templateObject.User,
	// 				Body:    l10n.GetL10n(channelInfo.GetChannelLang()).CantMuteModerator})
	// 			return
	// 		}

	// 	}
	// }

	// if template.StringRandomizer.Enabled == true {
	// 	if len(template.StringRandomizer.Strings) == 0 {
	// 		commandValues := strings.Split(chatCommand.Body, ",")
	// 		if len(commandValues) != 0 {
	// 			templateObject.RandomString = strings.TrimSpace(commandValues[rand.Intn(len(commandValues)-1)])

	// 		}
	// 	} else {
	// 		templateObject.RandomString = strings.TrimSpace(template.StringRandomizer.Strings[rand.Intn(len(template.StringRandomizer.Strings)-1)])
	// 	}
	// }

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
	log.Println(templateObject.PreventRedirect)
	log.Println(templateObject.PreventDebounce)
	//if chatCommand.Body != "" && !(template.StringRandomizer.Enabled == true && len(template.StringRandomizer.Strings) == 0) && template.PreventRedirect == false {
	if templateObject.PreventRedirect == false && chatCommand.Body != "" {
		if strings.HasPrefix(chatCommand.Body, "@") {
			redirectTo = chatCommand.Body[1:]
		} else {
			redirectTo = chatCommand.Body

		}
	}
	redirectTo = strings.Replace(redirectTo, "@", " @", -1)
	//if template.PreventDebounce == true {
	if templateObject.PreventDebounce == true {
		ircClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    templateObject.User,
			Body:    html.UnescapeString(message)})
		return
	}

	ircClient.SendDebounced(models.OutgoingDebouncedMessage{
		Message: models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    templateObject.User,
			Body:    html.UnescapeString(message)},
		Command:    template.AliasTo,
		RedirectTo: redirectTo})

}
