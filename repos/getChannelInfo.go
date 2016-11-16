package repos

import (
	"log"

	"github.com/khades/servbot/models"
)

// GetChannelInfo gets channel info
func GetChannelInfo(channel *string) (*models.ChannelInfo, error) {
	item, found := channelInfoRepositoryObject.dataArray[*channel]
	if found {
		return item, nil
	}
	var dbObject = &models.ChannelInfo{}
	error := Db.C("channelInfo").Find(models.ChannelSelector{Channel: *channel}).One(dbObject)
	if error != nil {
		log.Println("Error ", error)
		return nil, error
	}
	log.Println(dbObject)
	channelInfoRepositoryObject.dataArray[*channel] = dbObject
	return dbObject, error
}
