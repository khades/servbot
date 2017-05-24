package repos

import (
	"log"

	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

func PushVkGroupInfo(channelID *string, vkGroupInfo *models.VkGroupInfo) {
	log.Println("pushing info")
	log.Println(*channelID)
	log.Println(*vkGroupInfo)
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.VkGroupInfo = *vkGroupInfo
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, VkGroupInfo: *vkGroupInfo})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"vkgroupinfo": *vkGroupInfo}})
}

func GetVKEnabledChannels() (*[]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := Db.C(channelInfoCollection).Find(bson.M{"vkgroupinfo.groupname": bson.M{"$exists": true, "$ne": ""}}).All(&result)
	return &result, error
}
