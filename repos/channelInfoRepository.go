package repos

import "github.com/khades/servbot/models"

type channelInfoRepository struct {
	dataArray map[string]*models.ChannelInfo
}

func (c channelInfoRepository) forceCreateObject(channelID string, object *models.ChannelInfo) {
	c.dataArray[channelID] = object
}

var channelInfoRepositoryObject = channelInfoRepository{make(map[string]*models.ChannelInfo)}
