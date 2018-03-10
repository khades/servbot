package repos

import (
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

var channelInfoCollection = "channelInfo"

// GetChannelInfo gets channel info, and stores copy of that object in memory
func GetChannelInfo(channelID *string) (*models.ChannelInfo, error) {
	item, found := channelInfoRepositoryObject.dataArray[*channelID]
	if found {
		return item, nil
	}
	var dbObject = &models.ChannelInfo{}
	error := db.C(channelInfoCollection).Find(models.ChannelSelector{ChannelID: *channelID}).One(dbObject)
	if error != nil {
		log.Println("Error ", error)
		return nil, error
	}
	channelInfoRepositoryObject.dataArray[*channelID] = dbObject
	return dbObject, error
}

// GetModChannels returns list where specified user is moderator
func GetModChannels(userID *string) ([]models.ChannelWithID, error) {
	var result []models.ChannelWithID
	error := db.C(channelInfoCollection).Find(
		bson.M{"$or": []bson.M{
			bson.M{"mods": *userID},
			bson.M{"channelid": *userID}}}).All(&result)
	return result, error
}

// PreprocessChannels forces creation of empty channelInfo object in database, meanwhile preupdating channel names from IDs and writes everything to config
func PreprocessChannels() error {
	users, error := getTwitchUsersByID(Config.ChannelIDs)
	if error != nil {
		return error
	}
	channels := []string{}
	for _, user := range users {
		channels = append(channels, strings.ToLower(user.DisplayName))
		userIDCacheObject.Set("username-"+strings.ToLower(user.DisplayName), user.ID, 0)

		db.C(channelInfoCollection).Upsert(bson.M{"channelid": user.ID}, bson.M{
			"$set": bson.M{"channel": strings.ToLower(user.DisplayName)}})
	}
	Config.Channels = channels
	return nil
}

// PushStreamStatus updates stream status (start of stream, topic of stream)
func PushStreamStatus(channelID *string, streamStatus *models.StreamStatus) {
	channelInfo, _ := GetChannelInfo(channelID)
	prevStatus := models.StreamStatus{}
	currentStatus := models.StreamStatus{}
	if channelInfo != nil {
		prevStatus = channelInfo.StreamStatus
	}
	if prevStatus.Online == false && streamStatus.Online == true {

		currentStatus = models.StreamStatus{Online: true,
			Game:           streamStatus.Game,
			Title:          streamStatus.Title,
			Start:          streamStatus.Start,
			LastOnlineTime: time.Now(),
			Viewers:        streamStatus.Viewers,
			GamesHistory: []models.StreamStatusGameHistory{
				models.StreamStatusGameHistory{
					Game:  streamStatus.Game,
					Start: time.Now()}}}
	}
	if prevStatus.Online == true && streamStatus.Online == true {
		if prevStatus.Game == streamStatus.Game {
			currentStatus = models.StreamStatus{Online: true,
				Game:           streamStatus.Game,
				Title:          streamStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: time.Now(),
				Viewers:        streamStatus.Viewers,
				GamesHistory:   prevStatus.GamesHistory}
		} else {
			currentStatus = models.StreamStatus{Online: true,
				Game:           streamStatus.Game,
				Title:          streamStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: time.Now(),
				Viewers:        streamStatus.Viewers,
				GamesHistory: append(prevStatus.GamesHistory, models.StreamStatusGameHistory{
					Game:  streamStatus.Game,
					Start: time.Now()})}
		}
	}

	if streamStatus.Online == false {
		if prevStatus.Online == true && time.Since(prevStatus.LastOnlineTime).Seconds() < 120 {
			currentStatus = models.StreamStatus{Online: true,
				Game:           prevStatus.Game,
				Title:          prevStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: prevStatus.LastOnlineTime,
				Viewers:        prevStatus.Viewers,
				GamesHistory:   prevStatus.GamesHistory}
		} else {
			currentStatus = models.StreamStatus{
				Online:       false,
				Game:         "",
				Title:        "",
				Viewers:      0,
				GamesHistory: []models.StreamStatusGameHistory{}}
		}
	}
	if channelInfo != nil {
		channelInfo.StreamStatus = currentStatus
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, StreamStatus: currentStatus})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"streamstatus": currentStatus}})
}

// PushMods updates list of mods on channel
func PushMods(channelID *string, mods []string) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.Mods = mods
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, Mods: mods})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"mods": mods}})
}

// SetSubdayIsActive sets channelinfo flag "subdayisactive" to true
func SetSubdayIsActive(channelID *string, isActive bool) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.SubdayIsActive= isActive
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, SubdayIsActive: isActive})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subdayisactive": isActive}})
}

// PushCommandsForChannel splits all templates into offline and online arrays of channelInfo object
func PushCommandsForChannel(channelID *string) {
	var commandsList []string
	channelInfo, channelError := GetChannelInfo(channelID)
	dbCommands, commandsError := GetChannelActiveTemplates(channelID)
	var offlineCommands []string
	var onlineCommands []string
	if commandsError == nil {
		for _, value := range dbCommands {
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
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{
		"onlinecommands":  onlineCommands,
		"offlinecommands": offlineCommands}})
}
