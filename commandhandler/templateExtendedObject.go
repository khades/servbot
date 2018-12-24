package commandhandler

import (
	"fmt"
	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/followers"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/songRequest"
	"github.com/khades/servbot/subday"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type templateExtendedObject struct {
	channelInfo.ChannelInfo
	UserID string
	User   string
	IsMod              bool
	IsSub              bool
	CommandBody        string
	CommandBodyIsEmpty bool
	PreventDebounce    bool
	PreventRedirect    bool
	twitchIRCClient         *bot.TwitchIRCClient
	followersService *followers.Service
	subdayService *subday.Service
	songRequestService *songRequest.Service
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
	channelInfo.PreventRedirect = true
	isFollower, dur := channelInfo.followersService.IsFoller(&channelInfo.ChannelID, &channelInfo.UserID)
	if isFollower == true {
		return followerDuration{true, l10n.HumanizeDuration(time.Now().Sub(dur), channelInfo.Lang)}

	}
	return followerDuration{false, ""}
}

func (channelInfo templateExtendedObject) CurrentSong() channelInfo.CurrentSong {
	return channelInfo.songRequestService.GetLast(&channelInfo.ChannelID, channelInfo.Lang)
}

func (channelInfo *templateExtendedObject) SkipCurrentSong() string {
	channelInfo.PreventRedirect = true
	if channelInfo.IsMod == false {
		return l10n.GetL10n(channelInfo.Lang).SongRequestNotAModerator
	}
	channelInfo.PreventDebounce = true
	songrequest := channelInfo.songRequestService.GetLast(&channelInfo.ChannelID, channelInfo.Lang)
	if songrequest.IsPlaying == false {
		return l10n.GetL10n(channelInfo.GetChannelLang()).SongRequestNoRequests
	}
	channelInfo.songRequestService.Pull(&channelInfo.ChannelID, &songrequest.ID)
	return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulled, songrequest.Title)
}

func (channelInfo *templateExtendedObject) Random() string {
	lowerLimit := 0
	upperLimit := 100
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true
	if channelInfo.IsMod == false && channelInfo.StreamStatus.Online == true {
		return ""
	}
	return strconv.Itoa(lowerLimit + rand.Intn(upperLimit-lowerLimit+1))
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
		channelInfo.subdayService.CloseAnyActive(&channelInfo.ChannelID, &channelInfo.User, &channelInfo.UserID)
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
	banDuration := rand.Intn(length + 1)
	if banDuration == 0 {
		return banmeResult{
			Banned: false}
	}

	channelInfo.twitchIRCClient.SendPublic(&models.OutgoingMessage{
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
	if strings.TrimSpace(channelInfo.CommandBody) == "" {
		return ""
	}
	result := channelInfo.songRequestService.Add(&channelInfo.User, channelInfo.IsSub, &channelInfo.UserID, &channelInfo.ChannelID, &channelInfo.CommandBody)
	if result.Success == true {
		result.LengthStr = l10n.HumanizeDurationFull(result.Length, channelInfo.Lang, true)
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
	if result.TooShort == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestTooShort, result.Title)
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
	channelInfo.songRequestService.SetVolume(&channelInfo.ChannelID, int(volume))
	return l10n.GetL10n(channelInfo.Lang).VolumeChangeSuccess
}

func (channelInfo *templateExtendedObject) PullSongRequest() string {
	channelInfo.PreventDebounce = true
	channelInfo.PreventRedirect = true

	pulledVideo, pulled := channelInfo.songRequestService.PullLastUser(&channelInfo.ChannelID, &channelInfo.UserID)
	if pulled == true {
		return fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulled, pulledVideo.Title)
	}
	return l10n.GetL10n(channelInfo.Lang).SongRequestNoRequests

}
