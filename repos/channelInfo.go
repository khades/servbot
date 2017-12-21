package repos

import (
	"log"

	"gopkg.in/mgo.v2/bson"

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
	channelInfoRepositoryObject.dataArray[*channelID] = dbObject
	return dbObject, error
}

func GetModChannels(userID *string) (*[]models.ChannelWithID, error) {
	var result []models.ChannelWithID
	error := Db.C(channelInfoCollection).Find(
		bson.M{"$or": []bson.M{
			bson.M{"mods": *userID},
			bson.M{"channelid": *userID}}}).All(&result)
	return &result, error
}

func GetDubTrackEnabledChannels() (*[]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := Db.C(channelInfoCollection).Find(bson.M{"dubtrack.id": bson.M{"$exists": true, "$ne": ""}}).All(&result)
	return &result, error
}
