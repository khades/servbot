package repos

import (
	"log"
	"time"
)

// SetIfSubToChannel is cache for checking if user subbed to channel
func SetIfSubToChannel(userID *string, channelID *string, isSubbed *bool) {
	log.Printf("Setting cached value for user %s on channel %s: %s \n", *userID, *channelID, *isSubbed)
	duration := 2 * time.Hour
	if *isSubbed == false {
		duration = 5 * time.Minute
	}
	cacheObject.Set(*userID+*channelID+"isSub", *isSubbed, duration)
}
