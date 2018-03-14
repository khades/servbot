package models

import (
	"github.com/khades/servbot/l10n"
	"fmt"
	"math"
	"sort"
	"time"
	"unicode/utf8"
)

// StreamStatusGameHistory struct describes game on stream and its starting time
type StreamStatusGameHistory struct {
	Game   string    `bson:",omitempty"`
	GameID string    `bson:",omitempty"`
	Start  time.Time `bson:",omitempty"`
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

// StreamStatus Describes info about stream, when started, what game and title is, and if it is online
type StreamStatus struct {
	Online         bool
	Game           string       `bson:",omitempty"`
	GameID         string       `bson:",omitempty"`
	Title          string       `bson:",omitempty"`
	Start          time.Time    `bson:",omitempty"`
	LastOnlineTime time.Time    `bson:",omitempty"`
	Viewers        int          `bson:",omitempty"`
	GamesHistory   GamesHistory `bson:",omitempty"`
}
