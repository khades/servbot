package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

func GetChannelUsers(channelID *string, pattern *string) (*[]models.ChannelUsers, error) {
	var channelUsers []models.ChannelUsers
	if *pattern == "" {
		error := Db.C("messageLogs").Find(models.ChannelSelector{ChannelID: *channelID}).Sort("messages.date").Limit(100).All(&channelUsers)
		return &channelUsers, error

	}
	error := Db.C("messageLogs").Find(bson.M{
		"channelid": *channelID,
		"knownnicknames": bson.M{
			"$regex":   *pattern,
			"$options": "i"}}).Sort("messages.date").Limit(100).All(&channelUsers)
	return &channelUsers, error
}
