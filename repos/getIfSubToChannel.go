package repos

// GetIfSubToChannel is cache for checking if user subbed to channel
func GetIfSubToChannel(user *string, channel *string) (bool, bool) {
	value, found := cacheObject.Get(*user + *channel + "isSub")
	if found {
		return value.(bool), true

	}
	return false, false
}
