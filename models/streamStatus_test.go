package models

import (

	"sort"
	"testing"
	"time"
)

func TestStreamStatusOrdering(t *testing.T) {
	games := GamesHistory{
		StreamStatusGameHistory{Game: "2",
			Start: time.Now().Add(-4 * time.Minute)},
		StreamStatusGameHistory{Game: "3",
			Start: time.Now().Add(-3 * time.Minute)},
		StreamStatusGameHistory{Game: "1",
			Start: time.Now().Add(-5 * time.Minute)},

		StreamStatusGameHistory{Game: "4",
			Start: time.Now().Add(-2 * time.Minute)}}

	sort.Sort(sort.Reverse(games))
}
