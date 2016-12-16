package repos

import "github.com/khades/servbot/models"

func GetDubTrackEnabledChannels() []*models.ChannelInfo {
	result := []*models.ChannelInfo{}
	for _, value := range Config.Channels {
		value, error := GetChannelInfo(&value)
		if error == nil && value.DubTrack.ID != "" {
			result = append(result, value)
		}
	}
	return result
}
