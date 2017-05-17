package models

import (
	"fmt"
	"math"
	"sort"
	"time"
	"unicode/utf8"
)

type StreamStatusGameHistory struct {
	Game  string    `bson:",omitempty"`
	Start time.Time `bson:",omitempty"`
}
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

func (history GamesHistory) ReturnHistory() string {
	sort.Sort(sort.Reverse(history))
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
		minutes := float64(int(duration.Minutes() - math.Floor(duration.Minutes()/60)*60))
		hours := float64(int(duration.Hours()))
		stringDuration := fmt.Sprintf("%dh%dm]", int(hours), int(minutes))
		if minutes < 10 {
			stringDuration = fmt.Sprintf("%dh0%dm]", int(hours), int(minutes))
		}
		if stringHistory == "" {
			newStringHistory = "NOW" history[index].Game + " [" + stringDuration
		} else {
			newStringHistory = history[index].Game + " [" + stringDuration + " -> " + stringHistory
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
	Title          string       `bson:",omitempty"`
	Start          time.Time    `bson:",omitempty"`
	LastOnlineTime time.Time    `bson:",omitempty"`
	Viewers        int          `bson:",omitempty"`
	GamesHistory   GamesHistory `bson:",omitempty"`
}
