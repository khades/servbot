package streamStatus

import (
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/gameResolve"
	"github.com/khades/servbot/twitchAPIClient"
	"github.com/sirupsen/logrus"
)

type Service struct {
	config             *config.Config
	channelInfoService *channelInfo.Service
	gameResolveService *gameResolve.Service
	twitchAPIClient    *twitchAPIClient.TwitchAPIClient
}

func (service *Service) UpdateFromTwitch() error {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "twitchGames",
		"action":  "UpdateFromTwitch"})
	streams := make(map[string]channelInfo.StreamStatus)
	userIDs := service.config.ChannelIDs

	for _, channel := range userIDs {
		streams[channel] = channelInfo.StreamStatus{
			Online: false}
	}
	streamstatuses, error := service.twitchAPIClient.GetStreamStatuses()
	logger.Printf("Statuses: %+v", streamstatuses)
	if error != nil {
		return error
	}

	for _, status := range streamstatuses {
		streams[status.UserID] = channelInfo.StreamStatus{
			Online:  true,
			GameID:  status.GameID,
			Title:   status.Title,
			Start:   status.CreatedAt,
			Viewers: status.ViewerCount}

	}
	for channel, status := range streams {
		service.Push(&channel, &status)
	}
	return nil
}

// Push updates stream status (start of stream, topic of stream)
func (service *Service) Push(channelID *string, streamStatus *channelInfo.StreamStatus) {
	channelInfoData, _ := service.channelInfoService.GetChannelInfo(channelID)
	prevStatus := channelInfo.StreamStatus{}
	currentStatus := channelInfo.StreamStatus{}
	if channelInfoData != nil {
		prevStatus = channelInfoData.StreamStatus
	}
	game := ""
	if streamStatus.Online == true {
		game, _ = service.gameResolveService.Get(&streamStatus.GameID)

	}

	if prevStatus.Online == false && streamStatus.Online == true {

		currentStatus = channelInfo.StreamStatus{Online: true,
			Game:           game,
			GameID:         streamStatus.GameID,
			Title:          streamStatus.Title,
			Start:          streamStatus.Start,
			LastOnlineTime: time.Now(),
			Viewers:        streamStatus.Viewers,
			GamesHistory: []channelInfo.StreamStatusGameHistory{
				channelInfo.StreamStatusGameHistory{
					Game:   game,
					GameID: streamStatus.GameID,
					Start:  streamStatus.Start}}}
	}
	if prevStatus.Online == true && streamStatus.Online == true {
		if prevStatus.GameID == streamStatus.GameID {
			currentStatus = channelInfo.StreamStatus{Online: true,
				Game:           game,
				GameID:         streamStatus.GameID,
				Title:          streamStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: time.Now(),
				Viewers:        streamStatus.Viewers,
				GamesHistory:   prevStatus.GamesHistory}
		} else {
			currentStatus = channelInfo.StreamStatus{Online: true,
				Game:           game,
				GameID:         streamStatus.GameID,
				Title:          streamStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: time.Now(),
				Viewers:        streamStatus.Viewers,
				GamesHistory: append(prevStatus.GamesHistory, channelInfo.StreamStatusGameHistory{
					Game:   game,
					GameID: streamStatus.GameID,
					Start:  time.Now()})}
		}
	}

	if streamStatus.Online == false {
		if prevStatus.Online == true && time.Since(prevStatus.LastOnlineTime).Seconds() < 120 {
			currentStatus = channelInfo.StreamStatus{Online: true,
				Game:           game,
				GameID:         prevStatus.GameID,
				Title:          prevStatus.Title,
				Start:          prevStatus.Start,
				LastOnlineTime: prevStatus.LastOnlineTime,
				Viewers:        prevStatus.Viewers,
				GamesHistory:   prevStatus.GamesHistory}
		} else {
			currentStatus = channelInfo.StreamStatus{
				Online:       false,
				Game:         "",
				Title:        "",
				Viewers:      0,
				GamesHistory: []channelInfo.StreamStatusGameHistory{}}
		}
	}
	service.channelInfoService.SetStreamStatus(channelID, &currentStatus)

}
