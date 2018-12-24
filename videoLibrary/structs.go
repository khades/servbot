package videoLibrary

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

type TagRecord struct {
	Tag    string `json:"tag"`
	User   string `json:"user"`
	UserID string `json:"userID"`
}

type SongRequestLibraryResponse struct {
	Item *SongRequestLibraryItem 
	VideoDoesntExist bool
	InternalError bool
	VideoID    string
}