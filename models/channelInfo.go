package models

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ChannelInfo describes all information about channel
type ChannelInfo struct {
	Enabled             bool         `json:"enabled"`
	Lang                string       `json:"lang"`
	ChannelID           string       `json:"channelId"`
	Channel             string       `json:"channel"`
	StreamStatus        StreamStatus `json:"streamStatus"`
	TwitchDJ            TwitchDJ     `json:"twitchDJ"`
	VkGroupInfo         VkGroupInfo  `json:"vkGroupInfo"`
	Mods                []string     `json:"mods"`
	Commands            []string     `json:"commands"`
	ExtendedBansLogging bool         `json:"extendedBansLogging"`
	SubTrain            SubTrain     `json:"subTrain"`
	SubdayIsActive      bool         `json:"subdayIsActive"`
	SongRequest         CurrentSong  `json:"songRequest"`
}

// CurrentSong struct describes current state of songrequest on channel
type CurrentSong struct {
	IsPlaying bool   `json:"isPlaying"`
	Title     string `json:"title"`
	User      string `json:"user"`
	Link      string `json:"link"`
	Duration  string `json:"duration"`
	Volume    int    `json:"volume"`
	Count     int    `json:"count"`
}

// GetCommands Helper Command for mustashe
func (channelInfo ChannelInfo) GetCommands() string {
	return "!" + strings.Join(channelInfo.Commands, ", !")
}

// GetIfUserIsMod checks if user exist in internal mod array
func (channelInfo ChannelInfo) GetIfUserIsMod(userID *string) bool {
	isMod := false
	if *userID == "40635840" || channelInfo.ChannelID == *userID {
		return true
	}

	for _, value := range channelInfo.Mods {
		if value == *userID {
			isMod = true
			break
		}
	}
	return isMod
}

// GetChannelLang returns current channel language setting
func (channelInfo ChannelInfo) GetChannelLang() string {
	if channelInfo.Lang == "" {
		return "en"
	}
	return channelInfo.Lang
}

// GetGamesHistory returns formatted game history on channel
func (channelInfo ChannelInfo) GetGamesHistory() string {
	return channelInfo.StreamStatus.GamesHistory.ReturnHistory(channelInfo.GetChannelLang())
}

// GetStreamDuration Helper Command for time for mustashe
func (channelInfo ChannelInfo) GetStreamDuration() string {
	lang := channelInfo.GetChannelLang()
	if !channelInfo.StreamStatus.Online {
		return ""
	}
	minutePrefix := "минут"

	hourPrefix := "часов"
	if lang == "en" {
		minutePrefix = "minutes"
		hourPrefix = "hours"
	}
	duration := time.Now().Sub(channelInfo.StreamStatus.Start)
	minutes := float64(int(duration.Minutes() - math.Floor(duration.Minutes()/60)*60))
	hours := float64(int(duration.Hours()))
	if lang == "en" {
		if minutes == 1 {
			minutePrefix = "minute"
		}
		if hours == 1 {
			hourPrefix = "minute"
		}
	}
	if lang == "ru" && math.Floor(minutes/10) != 1 {
		switch int(minutes - math.Floor(minutes/10)*10) {
		case 1:
			minutePrefix = "минуту"

			break
		case 2:
		case 3:
		case 4:
			minutePrefix = "минуты"
			minutePrefix = "minutes"

		}
	}

	if lang == "ru" && int(math.Floor(hours/10)) != 1 {
		switch int(hours - math.Floor(hours/10)*10) {
		case 1:
			hourPrefix = "час"
			break
		case 2:
		case 3:
		case 4:
			hourPrefix = "часа"

		}
	}
	if int(minutes) == 0 {
		return fmt.Sprintf("%d %s", int(hours), hourPrefix)

	}
	if int(hours) == 0 {
		return fmt.Sprintf("%d %s", int(minutes), minutePrefix)
	}
	return fmt.Sprintf("%d %s %d %s", int(hours), hourPrefix, int(minutes), minutePrefix)

}
