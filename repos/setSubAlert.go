package repos

import (
	"time"

	"github.com/khades/servbot/models"
)

// SetSubAlert updates stream status (start of stream, topic of stream)
func SetSubAlert(user *string, subAlertInfo *models.SubAlertInfo) {
	ResubTemplateCache.Drop(&subAlertInfo.Channel)
	Db.C("subAlert").Upsert(models.ChannelSelector{Channel: subAlertInfo.Channel}, *subAlertInfo)
	Db.C("subAlertHistory").Insert(models.SubAlertHistory{*subAlertInfo, *user, time.Now()})
}
