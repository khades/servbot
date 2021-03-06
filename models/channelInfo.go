package models

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ChannelInfo describes all information about channel
type ChannelInfo struct {
	ChannelID       string       `json:"channelId"`
	Channel         string       `json:"channel"`
	StreamStatus    StreamStatus `json:"streamStatus"`
	Banme           Banme
	TwitchDJ        TwitchDJ    `json:"twitchDJ"`
	DubTrack        DubTrack    `json:"dubTrack"`
	VkGroupInfo     VkGroupInfo `json:"vkGroupInfo"`
	Mods            []string    `json:"mods"`
	OfflineCommands []string    `json:"offlinecommands"`
	OnlineCommands  []string    `json:"onlinecommands"`
	SubTrain        SubTrain    `json:"subTrain"`
	SubdayIsActive  bool        `json:"subdayIsActive"`
}
type ChannelInfoForTemplate struct {
	ChannelInfo
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

// GetCommands Helper Command for mustashe
func (channelInfo ChannelInfo) GetCommands() string {
	if channelInfo.StreamStatus.Online == true {
		return "!" + strings.Join(channelInfo.OnlineCommands, ", !")

	}
	return "!" + strings.Join(channelInfo.OfflineCommands, ", !")
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

func (chnanelInfo ChannelInfo) GetGamesHistory() string {
	return chnanelInfo.StreamStatus.GamesHistory.ReturnHistory()
}

// GetStreamDuration Helper Command for time for mustashe
func (channelInfo ChannelInfo) GetStreamDuration() string {

	if !channelInfo.StreamStatus.Online {
		return ""
	}
	minutePrefix := "минут"
	hourPrefix := "часов"
	duration := time.Now().Sub(channelInfo.StreamStatus.Start)
	minutes := float64(int(duration.Minutes() - math.Floor(duration.Minutes()/60)*60))
	hours := float64(int(duration.Hours()))
	if math.Floor(minutes/10) != 1 {
		switch int(minutes - math.Floor(minutes/10)*10) {
		case 1:
			minutePrefix = "минуту"
			break
		case 2:
		case 3:
		case 4:
			minutePrefix = "минуты"
		}
	}

	if int(math.Floor(hours/10)) != 1 {
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
