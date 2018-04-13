package repos

import (
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
)

var channelInfoCollection = "channelInfo"
var timesCalled int64

// setChannelName sets channel name after processing
func setChannelName(channelID *string, channel string) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.Channel = channel
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, Channel: channel})
	}
	db.C(channelInfoCollection).Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"channel": channel}})

}

// EnableChannel sets enabled flag of channel to true
func EnableChannel(channelID *string) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.Enabled = true
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, Enabled: true})
	}
	db.C(channelInfoCollection).Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"enabled": true}})

}

// DisableChannel sets enabled flag of channel to false
func DisableChannel(channelID *string) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.Enabled = false
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, Enabled: false})
	}
	db.C(channelInfoCollection).Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"enabled": false}})
}

//GetActiveChannels returns enabled channels on that chatbot instance
func GetActiveChannels() ([]models.ChannelInfo, error) {
	var results []models.ChannelInfo
	error := db.C(channelInfoCollection).Find(bson.M{"enabled": true}).All(&results)
	return results, error
}

// GetChannelInfo gets channel info, and stores copy of that object in memory
func GetChannelInfo(channelID *string) (*models.ChannelInfo, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "channelInfo",
		"action":  "GetChannelInfo"})
	timesCalled = timesCalled + 1
	logger.Debugf("Function was called %d times", timesCalled)
	item, found := channelInfoRepositoryObject.dataArray[*channelID]
	if found {
		return item, nil
	}
	var dbObject = &models.ChannelInfo{}
	error := db.C(channelInfoCollection).Find(models.ChannelSelector{ChannelID: *channelID}).One(dbObject)
	if error != nil {
		logger.Info("Error ", error)
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
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "channelInfo",
		"action":  "PreprocessChannels"})
	channels, channelError := GetActiveChannels()
	logger.Debugf("%d channels found", len(channels))
	channelIDList := []string{}
	if channelError != nil {
		logger.Debugf("Channels Error: %s", channelError.Error())
		return channelError
	}
	for _, channel := range channels {
		channelIDList = append(channelIDList, channel.ChannelID)
	}
	Config.ChannelIDs = channelIDList
	users, error := getTwitchUsersByID(channelIDList)
	if error != nil {
		logger.Debugf("Twitch Channels detection error: %s", error.Error())
		return error
	}
	channelsList := []string{}
	for _, user := range users {
		username := strings.ToLower(user.DisplayName)
		setChannelName(&user.ID, strings.ToLower(user.DisplayName))
		updateUserToUserIDFromChat(&user.ID, &username)
		channelsList = append(channelsList, user.DisplayName)
	}
	Config.Channels = channelsList
	return nil
}

// PushStreamStatus updates stream status (start of stream, topic of stream)
func pushStreamStatus(channelID *string, streamStatus *models.StreamStatus) {
	channelInfo, _ := GetChannelInfo(channelID)
	prevStatus := models.StreamStatus{}
	currentStatus := models.StreamStatus{}
	if channelInfo != nil {
		prevStatus = channelInfo.StreamStatus
	}
	if prevStatus.Online == false && streamStatus.Online == true {

		currentStatus = models.StreamStatus{Online: true,
			Game:           streamStatus.Game,
			GameID:         streamStatus.GameID,
			Title:          streamStatus.Title,
			Start:          streamStatus.Start,
			LastOnlineTime: time.Now(),
			Viewers:        streamStatus.Viewers,
			GamesHistory: []models.StreamStatusGameHistory{
				models.StreamStatusGameHistory{
					Game:   streamStatus.Game,
					GameID: streamStatus.GameID,
					Start:  streamStatus.Start}}}
	}
	if prevStatus.Online == true && streamStatus.Online == true {
		if prevStatus.GameID == streamStatus.GameID {
			currentStatus = models.StreamStatus{Online: true,
				Game:           streamStatus.Game,
				GameID:         streamStatus.GameID,
				Title:          streamStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: time.Now(),
				Viewers:        streamStatus.Viewers,
				GamesHistory:   prevStatus.GamesHistory}
		} else {
			currentStatus = models.StreamStatus{Online: true,
				Game:           streamStatus.Game,
				GameID:         streamStatus.GameID,
				Title:          streamStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: time.Now(),
				Viewers:        streamStatus.Viewers,
				GamesHistory: append(prevStatus.GamesHistory, models.StreamStatusGameHistory{
					Game:   streamStatus.Game,
					GameID: streamStatus.GameID,
					Start:  time.Now()})}
		}
	}

	if streamStatus.Online == false {
		if prevStatus.Online == true && time.Since(prevStatus.LastOnlineTime).Seconds() < 120 {
			currentStatus = models.StreamStatus{Online: true,
				Game:           prevStatus.Game,
				GameID:         prevStatus.GameID,
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

// SetChannelLang sets channel language
func SetChannelLang(channelID *string, lang *string) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.Lang = *lang
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, Lang: *lang})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"lang": *lang}})
}

// SetSubdayIsActive sets channelinfo flag "subdayisactive" to true
func SetSubdayIsActive(channelID *string, isActive bool) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.SubdayIsActive = isActive
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


	if commandsError == nil {
		for _, value := range dbCommands {
		
			commandsList = append(commandsList, value.CommandName)
		}
	}
	if channelError == nil {
		channelInfo.Commands =  commandsList
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{
			ChannelID:       *channelID,
			Commands:  commandsList,
		
		})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{
		"commands":  commandsList}})
}

func updateGamesByID(gameID *string, game *string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "twitchGames",
		"action":  "updateGamesByID"})
	logger.Debugf("Updating gameID %s with proper name \"%s\"", *gameID, *game)
	db.C(channelInfoCollection).UpdateAll(bson.M{"streamstatus.game": bson.M{"$exists": false}, "streamstatus.gameid": *gameID},
		bson.M{"$set": bson.M{"streamstatus.game": *game}})

	db.C(channelInfoCollection).UpdateAll(bson.M{"streamstatus.gameshistory.game": bson.M{"$exists": false}, "streamstatus.gameshistory.gameid": *gameID},
		bson.M{"$set": bson.M{"streamstatus.gameshistory.$.game": *game}})

	for _, channel := range channelInfoRepositoryObject.dataArray {
		if channel.StreamStatus.GameID == *gameID {
			channel.StreamStatus.Game = *game
		}
		for index, gameHistoryItem := range channel.StreamStatus.GamesHistory {
			logger.Debugf("Processing game history item with id %s", gameHistoryItem.GameID)

			if gameHistoryItem.GameID == *gameID {
				logger.Debugf("Found match for game history item with id %s", gameHistoryItem.GameID)
				channel.StreamStatus.GamesHistory[index].Game = *game
			}
		}
	}
}

// UpdateStreamStatuses gets all stream statuses for active channels and processes them before updating them into database
func UpdateStreamStatuses() error {
	streams := make(map[string]models.StreamStatus)
	userIDs := Config.ChannelIDs

	for _, channel := range userIDs {
		streams[channel] = models.StreamStatus{
			Online: false}
	}
	streamstatuses, error := getStreamStatuses()
	if error != nil {
		return error
	}

	for _, status := range streamstatuses {
		game, _ := getGameByID(&status.GameID)
		streams[status.UserID] = models.StreamStatus{
			Online:  true,
			GameID:  status.GameID,
			Game:    game,
			Title:   status.Title,
			Start:   status.CreatedAt,
			Viewers: status.ViewerCount}

	}
	for channel, status := range streams {
		pushStreamStatus(&channel, &status)
	}
	return nil
}
