package repos

import (
	"github.com/globalsign/mgo/bson"

	"github.com/khades/servbot/models"
)

// PushVkGroupInfo updates VK group info and last post content
func PushVkGroupInfo(channelID *string, vkGroupInfo *models.VkGroupInfo) {

	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.VkGroupInfo = *vkGroupInfo
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, VkGroupInfo: *vkGroupInfo})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"vkgroupinfo": *vkGroupInfo}})
}

// GetVKEnabledChannels returns list of channels, where VK group was configures
func GetVKEnabledChannels() ([]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := db.C(channelInfoCollection).Find(bson.M{"enabled":true,"vkgroupinfo.groupname": bson.M{"$exists": true, "$ne": ""}}).All(&result)
	return result, error
}
