package repos

import "github.com/khades/servbot/models"

func GetTwitchDJEnabledChannels() []*models.ChannelInfo {
	result := []*models.ChannelInfo{}
	for _, value := range Config.Channels {
		value, error := GetChannelInfo(&value)
		if error == nil && value.TwitchDJ.ID != "" {
			result = append(result, value)
		}
	}
	return result
}
