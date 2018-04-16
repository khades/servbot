package models

import "time"

// SongRequestLibraryItem describes already encountered songrequest
type SongRequestLibraryItem struct {
	VideoID    string        `json:"videoID"`
	Tags       []TagRecord   `json:"tags"`
	Length     time.Duration `json:"length"`
	Title      string        `json:"title"`
	ReviewedOn []string      `json:"reviewedOn"`
	LastCheck  time.Time     `json:"lastCheck"`
	Views      int64         `json:"views"`
	Likes      int64         `json:"likes"`
	Dislikes   int64         `json:"dislikes"`
}
