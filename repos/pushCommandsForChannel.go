package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// PushCommandsForChannel updates list of commands based on info from DB
func PushCommandsForChannel(channelID *string) {
	var commandsList []string
	channelInfo, channelError := GetChannelInfo(channelID)
	dbCommands, commandsError := GetChannelActiveTemplates(channelID)
	if commandsError == nil {
		for _, value := range *dbCommands {
			commandsList = append(commandsList, value.CommandName)
		}
	}
	if channelError == nil {
		channelInfo.Commands = commandsList
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, Commands: commandsList})
	}
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"commands": commandsList}})
}
