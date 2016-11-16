package repos

import "time"

// SetIfSubToChannel is cache for checking if user subbed to channel
func SetIfSubToChannel(user *string, channel *string, isSubbed *bool) {
	duration := 2 * time.Hour
	if *isSubbed == false {
		duration = 5 * time.Minute
	}
	cacheObject.Set(*user+*channel+"isSub", *isSubbed, duration)
}
