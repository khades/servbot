package repos

import "time"

// SetInfoSubToChannel is cache for checking if user subbed to channel
func SetInfoSubToChannel(user string, channel string, isSubbed bool) {
	duration := 2 * time.Hour
	if isSubbed == false {
		duration = 5 * time.Minute
	}
	cacheObject.Set(user+channel+"isSub", isSubbed, duration)
}
