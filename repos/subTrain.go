package repos

import (
	"github.com/globalsign/mgo/bson"
	"time"
	"github.com/khades/servbot/models"
)

// GetChannelsWithSubtrainNotification returns channels where subtrain notification should be shown
func GetChannelsWithSubtrainNotification() ([]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := db.C(channelInfoCollection).Find(
		bson.M{
			"subtrain.enabled": true,
			"subtrain.notificationshown" : false,
			"subtrain.currentstreak": bson.M{
				"$ne":0},
			"subtrain.notificationtime":bson.M{
				"$lt": time.Now()}}).All(&result)
	return result, error
}

// GetChannelsWithExpiredSubtrain returns channels where subtrain has expired
func GetChannelsWithExpiredSubtrain() ([]models.ChannelInfo, error) {
	result := []models.ChannelInfo{}
	error := db.C(channelInfoCollection).Find(
		bson.M{
			"subtrain.enabled": true,
			"subtrain.currentstreak": bson.M{
				"$ne":0},
			"subtrain.expirationtime":bson.M{
				"$lt": time.Now()}}).All(&result)
	return result, error
}

// PutChannelSubtrain upserts subtrain infromation for channel
func PutChannelSubtrain(channelID *string, subTrain *models.SubTrain) {
	channelInfo, _ := GetChannelInfo(channelID)
	if channelInfo != nil {
		channelInfo.SubTrain = *subTrain
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, SubTrain: *subTrain})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subtrain": *subTrain}})
}

// PutChannelSubtrainWeb upserts subtrain infromation for channel, unlike previous function, it tries to save current streak if possible
func PutChannelSubtrainWeb(channelID *string, subTrain *models.SubTrain) {
	channelInfo, _ := GetChannelInfo(channelID)
	localSubtrain := channelInfo.SubTrain
	if (subTrain.Enabled == true && localSubtrain.Enabled == true && localSubtrain.ExpirationLimit == subTrain.ExpirationLimit && localSubtrain.NotificationLimit == subTrain.NotificationLimit) {
		subTrain.ExpirationTime = localSubtrain.ExpirationTime
		subTrain.NotificationTime = localSubtrain.NotificationTime
		subTrain.CurrentStreak = localSubtrain.CurrentStreak
		subTrain.Users = localSubtrain.Users
		subTrain.NotificationShown = localSubtrain.NotificationShown
	}
	if channelInfo != nil {
		channelInfo.SubTrain = *subTrain
	} else {
		channelInfoRepositoryObject.forceCreateObject(*channelID, &models.ChannelInfo{ChannelID: *channelID, SubTrain: *subTrain})
	}
	db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"subtrain": *subTrain}})
}

// SetSubtrainNotificationShown sets "notificationshown" flag to true
func SetSubtrainNotificationShown(channelInfo *models.ChannelInfo) {
	subTrain := channelInfo.SubTrain
	subTrain.NotificationShown = true
	PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}	

// IncrementSubtrainCounterByChannelID is version of IncrementSubtrainCounter that gets channelInfo based on channelID
func IncrementSubtrainCounterByChannelID(channelID *string, user *string) {
	channelInfo, error := GetChannelInfo(channelID)
	if error == nil {
		IncrementSubtrainCounter(channelInfo, user)
		return
	} 

}

// IncrementSubtrainCounter increments specified channel subtrain information, also records subscriber username
func IncrementSubtrainCounter(channelInfo *models.ChannelInfo, user *string) {
	subTrain := channelInfo.SubTrain
	if subTrain.Enabled == false {
		return
	}
	subTrain.ExpirationTime = time.Now().Add(time.Second * time.Duration(subTrain.ExpirationLimit))
	subTrain.NotificationTime = time.Now().Add(time.Second * time.Duration(subTrain.NotificationLimit))
	subTrain.CurrentStreak = subTrain.CurrentStreak +1
	subTrain.Users = append(subTrain.Users, *user)
	PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}

// ResetSubtrainCounter resets current subtrain counters
func ResetSubtrainCounter(channelInfo *models.ChannelInfo) {
	subTrain := channelInfo.SubTrain
	subTrain.CurrentStreak = 0
	subTrain.NotificationShown = false
	subTrain.Users = []string{}
	PutChannelSubtrain(&channelInfo.ChannelID, &subTrain)
}