package repos

import (
//	"time"


)


// GetIfSubToChannel is cache for checking if user subbed to channel
func GetIfSubToChannel(userID *string, channelID *string) (bool, bool) {
	// value, found := cacheObject.Get(*userID + *channelID + "isSub")
	// if found {
	// 	return value.(bool), true
	// }
	return false, false
}

// SetIfSubToChannel is cache for checking if user subbed to channel
// func SetIfSubToChannel(userID *string, channelID *string, isSubbed *bool) {
// 	// duration := 2 * time.Hour
// 	// if *isSubbed == false {
// 	// 	duration = 5 * time.Minute
// 	// }
// 	// cacheObject.Set(*userID+*channelID+"isSub", *isSubbed, duration)
// }
