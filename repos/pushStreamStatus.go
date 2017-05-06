package repos

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

// PushStreamStatus updates stream status (start of stream, topic of stream)
func PushStreamStatus(channelID *string, streamStatus *models.StreamStatus) {
	channelInfo, _ := GetChannelInfo(channelID)
	prevStatus := models.StreamStatus{}
	currentStatus := models.StreamStatus{}
	if channelInfo != nil {
		prevStatus = channelInfo.StreamStatus
	}
	log.Println("Outputting prev status")
	log.Println(prevStatus)
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
		if prevStatus.Game != streamStatus.Game {
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
	Db.C(channelInfoCollection).Upsert(models.ChannelSelector{ChannelID: *channelID}, bson.M{"$set": bson.M{"streamstatus": currentStatus}})
}
