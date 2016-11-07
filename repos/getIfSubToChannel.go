package repos

// GetInfoSubToChannel is cache for checking if user subbed to channel
func GetInfoSubToChannel(key string, channel string) (bool, bool) {
	value, found := cacheObject.Get(key + channel + "isSub")
	return value.(bool), found
}
