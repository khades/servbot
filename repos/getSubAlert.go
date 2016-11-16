package repos

import "github.com/khades/servbot/models"

// GetSubAlert ha
func GetSubAlert(channel *string) (*models.SubAlertInfo, error) {
	var result models.SubAlertInfo
	error := Db.C("subAlert").Find(models.ChannelSelector{Channel: *channel}).One(&result)
	return &result, error
}
