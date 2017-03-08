package repos

import (
	"log"

	"github.com/khades/servbot/models"
)

// GetChannelInfo gets channel info
func GetChannelInfo(channelID *string) (*models.ChannelInfo, error) {
	item, found := channelInfoRepositoryObject.dataArray[*channelID]
	if found {
		return item, nil
	}
	var dbObject = &models.ChannelInfo{}
	error := Db.C(channelInfoCollection).Find(models.ChannelSelector{ChannelID: *channelID}).One(dbObject)
	if error != nil {
		log.Println("Error ", error)
		return nil, error
	}
	//log.Println(dbObject)
	channelInfoRepositoryObject.dataArray[*channelID] = dbObject
	return dbObject, error
}
