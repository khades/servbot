package videoLibraryAPI

import "github.com/khades/servbot/videoLibrary"

type videolibraryResponse struct {
	Count int                             `json:"count"`
	Items []videoLibrary.SongRequestLibraryItem `json:"items"`
}
