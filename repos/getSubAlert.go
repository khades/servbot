package repos

import "github.com/khades/servbot/models"

// GetSubAlert ha
func GetSubAlert(channelID *string) (*models.SubAlert, error) {
	var result models.SubAlert
	error := Db.C(subAlertCollection).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}

func GetSubAlertWithHistory(channelID *string) (*models.SubAlertWithHistory, error) {
	var result models.SubAlertWithHistory
	error := Db.C(subAlertCollection).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}
