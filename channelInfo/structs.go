package channelInfo

import (
	"fmt"
	"github.com/khades/servbot/l10n"
	"math"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// ChannelInfo describes all information about channel
type ChannelInfo struct {
	Enabled             bool         `json:"enabled"`
	Lang                string       `json:"lang"`
	ChannelID           string       `json:"channelID"`
	Channel             string       `json:"channel"`
	StreamStatus        StreamStatus `json:"streamStatus"`
	// TwitchDJ            TwitchDJ     `json:"twitchDJ"`
	VkGroupInfo         VkGroupInfo  `json:"vkGroupInfo"`
	Mods                []string     `json:"mods"`
	Commands            []string     `json:"commands"`
	ExtendedBansLogging bool         `json:"extendedBansLogging"`
	SubTrain            SubTrain     `json:"subTrain"`
	SubdayIsActive      bool         `json:"subdayIsActive"`
	CurrentSong        CurrentSong  `json:"currentSong"`
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


// StreamStatus Describes info about stream, when started, what game and title is, and if it is online
type StreamStatus struct {
	Online         bool         `json:"online"`
	Game           string       `bson:",omitempty" json:"game"`
	GameID         string       `bson:",omitempty" json:"gameID"`
	Title          string       `bson:",omitempty" json:"title"`
	Start          time.Time    `bson:",omitempty" json:"start"`
	LastOnlineTime time.Time    `bson:",omitempty" json:"lastOnlineTime"`
	Viewers        int          `bson:",omitempty" json:"viewers"`
	GamesHistory   GamesHistory `bson:",omitempty" json:"gamesHistory"`
}

// StreamStatusGameHistory struct describes game on stream and its starting time
type StreamStatusGameHistory struct {
	Game   string    `bson:",omitempty" json:"game"`
	GameID string    `bson:",omitempty" json:"gameID"`
	Start  time.Time `bson:",omitempty" json:"start"`
}

// GamesHistory is type alias to array of StreamStatusGameHistory, used to properly sorting history by date
type GamesHistory []StreamStatusGameHistory

func (history GamesHistory) Len() int {
	return len(history)
}

func (history GamesHistory) Less(i, j int) bool {
	return history[i].Start.Before(history[j].Start)
}

func (history GamesHistory) Swap(i, j int) {
	history[i], history[j] = history[j], history[i]
}

// ReturnHistory forms one-line human-readable history of games on running stream
func (history GamesHistory) ReturnHistory(lang string) string {
	sort.Sort(sort.Reverse(history))
	var mPrefix = "m"
	var hPrefix = "h"
	var nowPrefix = l10n.GetL10n(lang).GameHistoryActivePrefix
	if lang == "ru" {
		mPrefix = "м"
		hPrefix = "ч"

	}
	stringHistory := ""
	if len(history) == 0 {
		return ""
	}
	for index := range history {
		newStringHistory := ""

		duration := time.Second * 1
		if index == 0 {
			duration = time.Now().Sub(history[0].Start)
		} else {
			duration = history[index-1].Start.Sub(history[index].Start)
		}
		minutes := int(duration.Minutes() - math.Floor(duration.Minutes()/60)*60)
		hours := int(duration.Hours())
		stringDuration := "["
		if hours > 0 {
			stringDuration = fmt.Sprintf("[%d%s", hours, hPrefix)
		}

		if minutes < 10 {
			stringDuration = stringDuration + fmt.Sprintf("0%d%s]", minutes, mPrefix)
		} else {
			stringDuration = stringDuration + fmt.Sprintf("%d%s]", minutes, mPrefix)
		}

		if stringHistory == "" {
			newStringHistory = nowPrefix + " " + history[index].Game + " " + stringDuration
		} else {
			newStringHistory = history[index].Game + " " + stringDuration + " > " + stringHistory
		}

		if utf8.RuneCountInString(newStringHistory) > 400 {
			break
		} else {
			stringHistory = newStringHistory
		}
	}
	return stringHistory
}


// SubTrain descibes subtrain settings and state on specified channel
type SubTrain struct {
	Enabled              bool      `json:"enabled"`
	OnlyNewSubs          bool      `json:"onlyNewSubs"`
	ExpirationLimit      int       `json:"expirationLimit"`
	NotificationLimit    int       `json:"notificationLimit"`
	NotificationShown    bool      `json:"notificationShown"`
	ExpirationTime       time.Time `json:"expirationTime"`
	NotificationTime     time.Time `json:"notificationTime"`
	AppendTemplate       string    `json:"appendTemplate"`
	TimeoutTemplate      string    `json:"timeoutTemplate"`
	NotificationTemplate string    `json:"notificationTemplate"`
	CurrentStreak        int       `json:"currentStreak"`
	Users                []string  `json:"users"`
}

// VkGroupInfo describes information of vk group
type VkGroupInfo struct {
	GroupName       string `json:"groupName"`
	NotifyOnChange  bool   `json:"notifyOnChange"`
	LastMessageID   int    `json:"lastMessageID"`
	LastMessageURL  string `json:"lastMessageURL"`
	LastMessageBody string `json:"lastMessageBody"`
	LastMessageDate string `json:"lastMessageDate"`
}

// TwitchDJ describes information about twitchDJ service
type TwitchDJ struct {
	ID             string `json:"id"`
	Playing        bool   `json:"playing"`
	Track          string `json:"track"`
	NotifyOnChange bool   `json:"notifyOnChange"`
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
	ID        string `json:"id"`
}

type ChannelWithID struct {
	Channel   string `json:"channel"`
	ChannelID string `json:"channelID"`
}
