package channelInfo

import (
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/userResolve"
	"github.com/khades/servbot/utils"
	"github.com/sirupsen/logrus"
)

type Service struct {
	// Dependencies
	collection         *mgo.Collection
	config             *config.Config
	userResolveService *userResolve.Service

	// Own Fields
	dataArray map[string]*ChannelInfo
}

func (c *Service) forceCreateObject(channelID string, object *ChannelInfo) {
	c.dataArray[channelID] = object
}

//var с = Service{make(map[string]*models.ChannelInfo)}

// setChannelName sets channel name after processing
func (c *Service) setChannelName(channelID *string, channel string) {
	channelInfo, _ := c.Get(channelID)
	if channelInfo != nil {
		channelInfo.Channel = channel
	} else {
		c.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, Channel: channel})
	}
	c.collection.Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"channel": channel}})

}

// Enable sets enabled flag of channel to true
func (c *Service) Enable(channelID *string) {
	channelInfo, _ := c.Get(channelID)
	if channelInfo != nil {
		channelInfo.Enabled = true
	} else {
		c.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, Enabled: true})
	}
	c.collection.Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"enabled": true}})

}

// Disable sets enabled flag of channel to false
func (c *Service) Disable(channelID *string) {
	channelInfo, _ := c.Get(channelID)
	if channelInfo != nil {
		channelInfo.Enabled = false
	} else {
		c.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, Enabled: false})
	}
	c.collection.Upsert(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"enabled": false}})
}

// TODO: Fix to get if bot is mod on channel
func (c *Service) GetChannelsWithExtendedLogging() ([]ChannelInfo, error) {
	var results []ChannelInfo
	error := c.collection.Find(bson.M{"enabled": true, "extendedbanslogging": true}).All(&results)
	return results, error
}

//GetActiveChannels returns enabled channels on that chatbot instance
func (c *Service) GetActiveChannels() ([]ChannelInfo, error) {
	var results []ChannelInfo
	error := c.collection.Find(bson.M{"enabled": true}).All(&results)
	return results, error
}

func (service *Service) PutCurrentSong(channelID *string, currentSong *CurrentSong) {
	item, found := service.dataArray[*channelID]
	if found {
		item.CurrentSong = *currentSong
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, CurrentSong: *currentSong})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"currentsong": *currentSong}})
}

// c.Get gets channel info, and stores copy of that object in memory
func (c *Service) Get(channelID *string) (*ChannelInfo, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "channelInfo",
		"action":  "get"})
	//	logger.Debugf("Function was called %d times", timesCalled)
	logger.Debugf("Getting channel :%s", *channelID)
	item, found := c.dataArray[*channelID]
	if found {
		return item, nil
	}
	var dbObject = ChannelInfo{}
	error := c.collection.Find(utils.ChannelSelector{ChannelID: *channelID}).One(dbObject)
	if error != nil {
		logger.Debugf("Error %s", error.Error())
		return nil, error
	}
	c.dataArray[*channelID] = &dbObject
	return &dbObject, error
}

// GetModChannels returns list where specified user is moderator
func (c *Service) GetModChannels(userID *string) ([]ChannelWithID, error) {
	// Needs check if channel is enabled
	var result []ChannelWithID
	error := c.collection.Find(
		bson.M{"$or": []bson.M{
			bson.M{"mods": *userID},
			bson.M{"channelid": *userID}}}).All(&result)
	return result, error
}

func (c *Service) GetModChannelsForAdmin() ([]ChannelWithID, error) {
	var result []ChannelWithID
	error := c.collection.Find(
		bson.M{"enabled": true}).All(&result)
	return result, error
}

// PreprocessChannels forces creation of empty channelInfo object in database, meanwhile preupdating channel names from IDs and writes everything to config
func (c *Service) PreprocessChannels() error {
	logger := logrus.WithFields(logrus.Fields{
		"package": "channelInfo",
		"action":  "PreprocessChannels"})
	channels, channelError := c.GetActiveChannels()
	logger.Debugf("%d channels found", len(channels))
	channelIDList := []string{}
	if channelError != nil {
		logger.Debugf("Channels Error: %s", channelError.Error())
		return channelError
	}

	for index, channel := range channels {
		channelIDList = append(channelIDList, channel.ChannelID)
		c.dataArray[channel.ChannelID] = &channels[index]
	}

	users, error := c.userResolveService.GetUsernames(channelIDList)
	if error != nil {
		logger.Debugf("Twitch Channels detection error: %s", error.Error())
		return error
	}

	channelsList := []string{}
	for userID, user := range *users {
		username := strings.ToLower(user)
		c.setChannelName(&userID, strings.ToLower(user))
		c.userResolveService.Update(&userID, &username)
		channelsList = append(channelsList, user)
	}
	c.config.ChannelIDs = channelIDList
	c.config.Channels = channelsList
	return nil
}

func (service *Service) SetStreamStatus(channelID *string, streamStatus *StreamStatus) {
	item, found := service.dataArray[*channelID]
	if found {
		item.StreamStatus = *streamStatus
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, StreamStatus: *streamStatus})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"streamstatus": *streamStatus}})
}

// PushMods updates list of mods on channel
func (service *Service) PushMods(channelID *string, mods []string) {
	channelInfo, _ := service.Get(channelID)
	if channelInfo != nil {
		channelInfo.Mods = mods
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, Mods: mods})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"mods": mods}})
}

// SetChannelLang sets channel language
func (service *Service) SetChannelLang(channelID *string, lang *string) {
	channelInfo, _ := service.Get(channelID)
	if channelInfo != nil {
		channelInfo.Lang = *lang
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, Lang: *lang})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"lang": *lang}})
}

// SetSubdayIsActive sets channelinfo flag "subdayisactive" to true
func (service *Service) SetSubdayIsActive(channelID *string, isActive bool) {
	channelInfo, _ := service.Get(channelID)
	if channelInfo != nil {
		channelInfo.SubdayIsActive = isActive
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, SubdayIsActive: isActive})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subdayisactive": isActive}})
}

// SetCommandsForChannel sets active commands for channel
func (service *Service) SetCommandsForChannel(channelID *string, commandsList []string) {
	channelInfo, channelError := service.Get(channelID)

	if channelError == nil {
		channelInfo.Commands = commandsList
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{
			ChannelID: *channelID,
			Commands:  commandsList,
		})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{
		"commands": commandsList}})
}

// UpdateGamesByID updates gameIDs to full game names, when gameid->game resolution is done
func (service *Service) UpdateGamesByID(gameID *string, game *string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "channelInfo",
		"action":  "UpdateGamesByID"})
	logger.Debugf("Updating gameID %s with proper name \"%s\"", *gameID, *game)
	service.collection.UpdateAll(bson.M{"streamstatus.game": bson.M{"$exists": false}, "streamstatus.gameid": *gameID},
		bson.M{"$set": bson.M{"streamstatus.game": *game}})

	service.collection.UpdateAll(bson.M{"streamstatus.gameshistory.game": bson.M{"$exists": false}, "streamstatus.gameshistory.gameid": *gameID},
		bson.M{"$set": bson.M{"streamstatus.gameshistory.$.game": *game}})

	for _, channel := range service.dataArray {
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

//  GetChannelNameByID tries to fetch channelname for specified channelid from cache, falling back to channelInfo database
func (c *Service) GetChannelNameByID(channelID *string) (*string, error) {

	logger := logrus.WithFields(logrus.Fields{
		"package": "channelInfo",
		"action":  "GetChannelNameByID"})

	logger.Debugf("Looking for channel by id %s", *channelID)

	channel, channelError := c.Get(channelID)

	if channelError != nil {
		logger.Debugf("Looking for channelid %s in database failed: %s", *channelID, channelError.Error())

		return nil, channelError
	}

	userName := strings.ToLower(channel.Channel)
	logger.Debugf("ChannelID %s has name: %s", *channelID, strings.ToLower(channel.Channel))

	return &userName, nil
}

// SUBTRAINS METHODS
// GetChannelsWithSubtrainNotification returns channels where subtrain notification should be shown
func (service *Service) GetChannelsWithSubtrainNotification() ([]ChannelInfo, error) {
	result := []ChannelInfo{}
	error := service.collection.Find(
		bson.M{
			"enabled":                    true,
			"subtrain.enabled":           true,
			"subtrain.notificationshown": false,
			"subtrain.currentstreak": bson.M{
				"$ne": 0},
			"subtrain.notificationtime": bson.M{
				"$lt": time.Now()}}).All(&result)
	return result, error
}

// GetChannelsWithExpiredSubtrain returns channels where subtrain has expired
func (service *Service) GetChannelsWithExpiredSubtrain() ([]ChannelInfo, error) {
	result := []ChannelInfo{}
	error := service.collection.Find(
		bson.M{
			"enabled":          true,
			"subtrain.enabled": true,
			"subtrain.currentstreak": bson.M{
				"$ne": 0},
			"subtrain.expirationtime": bson.M{
				"$lt": time.Now()}}).All(&result)
	return result, error
}

// PutChannelSubtrain upserts subtrain infromation for channel
func (service *Service) PutChannelSubtrain(channelID *string, subTrain *SubTrain) {
	channelInfo, _ := service.Get(channelID)
	if channelInfo != nil {
		channelInfo.SubTrain = *subTrain
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, SubTrain: *subTrain})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subtrain": *subTrain}})
}

// PutChannelSubtrainWeb upserts subtrain infromation for channel, unlike previous function, it tries to save current streak if possible
func (service *Service) PutChannelSubtrainWeb(channelID *string, subTrain *SubTrain) {
	channelInfo, _ := service.Get(channelID)
	localSubtrain := channelInfo.SubTrain
	if subTrain.Enabled == true && localSubtrain.Enabled == true && localSubtrain.ExpirationLimit == subTrain.ExpirationLimit && localSubtrain.NotificationLimit == subTrain.NotificationLimit {
		subTrain.ExpirationTime = localSubtrain.ExpirationTime
		subTrain.NotificationTime = localSubtrain.NotificationTime
		subTrain.CurrentStreak = localSubtrain.CurrentStreak
		subTrain.Users = localSubtrain.Users
		subTrain.NotificationShown = localSubtrain.NotificationShown
	}
	if channelInfo != nil {
		channelInfo.SubTrain = *subTrain
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, SubTrain: *subTrain})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subtrain": *subTrain}})
}

// SetSubtrainNotificationShown sets "notificationshown" flag to true
func (service *Service) SetSubtrainNotificationShown(channelInfo *ChannelInfo) {
	subTrain := channelInfo.SubTrain
	subTrain.NotificationShown = true
	service.PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}

// IncrementSubtrainCounterByChannelID is version of IncrementSubtrainCounter that gets channelInfo based on channelID
func (service *Service) IncrementSubtrainCounterByChannelID(channelID *string, user *string) {
	channelInfo, error := service.Get(channelID)
	if error == nil {
		service.IncrementSubtrainCounter(channelInfo, user)
		return
	}

}

// IncrementSubtrainCounter increments specified channel subtrain information, also records subscriber username
func (service *Service) IncrementSubtrainCounter(channelInfo *ChannelInfo, user *string) {
	subTrain := channelInfo.SubTrain
	if subTrain.Enabled == false {
		return
	}
	subTrain.ExpirationTime = time.Now().Add(time.Second * time.Duration(subTrain.ExpirationLimit))
	subTrain.NotificationTime = time.Now().Add(time.Second * time.Duration(subTrain.NotificationLimit))
	subTrain.CurrentStreak = subTrain.CurrentStreak + 1
	subTrain.Users = append(subTrain.Users, *user)
	service.PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}

// ResetSubtrainCounter resets current subtrain counters
func (service *Service) ResetSubtrainCounter(channelInfo *ChannelInfo) {
	subTrain := channelInfo.SubTrain
	subTrain.CurrentStreak = 0
	subTrain.NotificationShown = false
	subTrain.Users = []string{}
	service.PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}

// VK!!
// PushVkGroupInfo updates VK group info and last post content
func (service *Service) PushVkGroupInfo(channelID *string, vkGroupInfo *VkGroupInfo) {

	channelInfo, _ := service.Get(channelID)
	if channelInfo != nil {
		channelInfo.VkGroupInfo = *vkGroupInfo
	} else {
		service.forceCreateObject(*channelID, &ChannelInfo{ChannelID: *channelID, VkGroupInfo: *vkGroupInfo})
	}
	service.collection.Upsert(utils.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"vkgroupinfo": *vkGroupInfo}})
}

// GetVKEnabledChannels returns list of channels, where VK group was configures
func (service *Service) GetVKEnabledChannels() ([]ChannelInfo, error) {
	result := []ChannelInfo{}
	error := service.collection.Find(bson.M{"enabled": true, "vkgroupinfo.groupname": bson.M{"$exists": true, "$ne": ""}}).All(&result)
	return result, error
}
