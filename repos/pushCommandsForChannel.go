package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// PushCommandsForChannel updates list of commands based on info from DB
func PushCommandsForChannel(channel *string) {
	var commandsList []string
	channelInfo, channelError := GetChannelInfo(channel)
	dbCommands, commandsError := GetChannelActiveTemplates(channel)
	if commandsError == nil {
		for _, value := range *dbCommands {
			commandsList = append(commandsList, value.CommandName)
		}
	}
	if channelError == nil {
		channelInfo.Commands = commandsList
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channel, &models.ChannelInfo{Channel: *channel, Commands: commandsList})
	}
	Db.C("channelInfo").Upsert(models.ChannelSelector{Channel: *channel}, bson.M{"$set": bson.M{"commands": commandsList}})
}
