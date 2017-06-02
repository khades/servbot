package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

func PushCommandsForChannel(channelID *string) {
	var commandsList []string
	channelInfo, channelError := GetChannelInfo(channelID)
	dbCommands, commandsError := GetChannelActiveTemplates(channelID)
	var offlineCommands []string
	var onlineCommands []string
	if commandsError == nil {
		for _, value := range *dbCommands {
			if value.ShowOffline == true {
				offlineCommands = append(offlineCommands, value.CommandName)
			}
			if value.ShowOnline == true {
				onlineCommands = append(onlineCommands, value.CommandName)
			}
			commandsList = append(commandsList, value.CommandName)
		}
	}
	if channelError == nil {
		channelInfo.OfflineCommands = offlineCommands
		channelInfo.OnlineCommands = onlineCommands
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{
			ChannelID:       *channelID,
			OnlineCommands:  onlineCommands,
			OfflineCommands: offlineCommands,
		})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{
		"onlinecommands":  onlineCommands,
		"offlinecommands": offlineCommands}})
}
