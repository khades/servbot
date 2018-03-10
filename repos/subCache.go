package repos

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var cacheObject = cache.New(15*time.Minute, 30*time.Second)

// GetIfSubToChannel is cache for checking if user subbed to channel
func GetIfSubToChannel(userID *string, channelID *string) (bool, bool) {
	value, found := cacheObject.Get(*userID + *channelID + "isSub")
	if found {
		return value.(bool), true
	}
	return false, false
}
