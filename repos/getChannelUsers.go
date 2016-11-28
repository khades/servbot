package repos

import "github.com/khades/servbot/models"

func GetChannelUsers(channel *string) (*[]models.ChannelUserSelector, error) {
	var channelUsers []models.ChannelUserSelector
	error := Db.C("messageLogs").Find(models.ChannelSelector{Channel: *channel}).All(&channelUsers)
	return &channelUsers, error
}
