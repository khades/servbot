package repos

import "github.com/khades/servbot/models"

func GetChannelUsers(channelID *string) (*[]models.ChannelUsers, error) {
	var channelUsers []models.ChannelUsers
	error := Db.C("messageLogs").Find(models.ChannelSelector{ChannelID: *channelID}).All(&channelUsers)
	return &channelUsers, error
}
