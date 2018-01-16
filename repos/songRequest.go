package repos

import (
	//"time"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var songRequestCollectionName = "songrequests"

func GetSongRequest(channelID *string) *models.ChannelSongRequest {
	songRequestInfo := models.ChannelSongRequest{Settings: models.ChannelSongRequestSettings{PlaylistLength: 30, MaxVideoLength: 300, MaxRequestsPerUser: 2}}
	Db.C(songRequestCollectionName).Find(
		bson.M{
			"channelid": *channelID}).One(&songRequestInfo)
	if songRequestInfo.Settings.PlaylistLength == 0 {
		songRequestInfo.Settings.PlaylistLength = 30

	}
	if songRequestInfo.Settings.MaxVideoLength == 0 {
		songRequestInfo.Settings.MaxVideoLength = 300

	}
	if songRequestInfo.Settings.MaxRequestsPerUser == 0 {
		songRequestInfo.Settings.MaxRequestsPerUser = 3

	}
	if songRequestInfo.Settings.VideoViewLimit == 0 {
		songRequestInfo.Settings.VideoViewLimit  = 2000

	}
	return &songRequestInfo
}

func PushSongRequest(channelID *string, request *models.SongRequest) {
	Db.C(songRequestCollectionName).Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$push": bson.M{"requests": *request}})
}

func PushSongRequestSettings(channelID *string, settings *models.ChannelSongRequestSettings) {
	Db.C(songRequestCollectionName).Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$set": bson.M{"settings": *settings}})
}

func AddSongRequest(user *string, userIsSub bool, userID *string, channelID *string, videoID *string) error {
	songRequestInfo := GetSongRequest(channelID)

	if len(songRequestInfo.Requests) >= songRequestInfo.Settings.PlaylistLength {
		return errors.New("Playlist is full")
	}
	songIndex := sort.Search(len(songRequestInfo.Requests), func(i int) bool { return songRequestInfo.Requests[i].VideoID == *videoID })
	if songIndex != len(songRequestInfo.Requests) {
		return errors.New("Song is already in playlist")
	}
	userRequestsCount := 0
	for _, request := range songRequestInfo.Requests {
		if request.UserID == *userID {
			userRequestsCount = userRequestsCount + 1
		}
	}
	if userRequestsCount >= songRequestInfo.Settings.MaxRequestsPerUser {
		return errors.New("Too many requests per user")
	}
	video, videoError := GetYoutubeVideoInfo(videoID)
	if videoError != nil {
		return videoError
	}
	if len(video.Items) == 0 {
		return errors.New("Nothing found")
	}
	duration, durationError := video.Items[0].ContentDetails.GetDuration()
	if durationError != nil {
		return durationError
	}
	if duration.Seconds() > float64(songRequestInfo.Settings.MaxVideoLength) {
		return errors.New("Video is too long")

	}
	if video.Items[0].Statistics.GetViewCount() < songRequestInfo.Settings.VideoViewLimit {
		return errors.New("Too little views on video, got "+string(video.Items[0].Statistics.ViewCount) )

	}
	songRequest := models.SongRequest{User: *user, UserID: *userID, Date: time.Now(), VideoID: *videoID, Length: *duration, Title: video.Items[0].Snippet.Title}
	PushSongRequest(channelID, &songRequest)
	return errors.New("Video " + video.Items[0].Snippet.Title + " , length: " + duration.String())
}

func PullSongRequest(channelID *string, videoID *string) {
	Db.C(songRequestCollectionName).Update(
		bson.M{
			"channelid": *channelID}, bson.M{"$pull": bson.M{"requests":  bson.M{"videoid":*videoID}}})
}

func PullUserSongRequest(channelID *string, videoID *string, userID *string) {
	Db.C(songRequestCollectionName).Update(
		bson.M{
			"channelid": *channelID}, bson.M{"$pull": bson.M{"requests":  bson.M{"userid:":*userID,"videoid":*videoID}}})
}

func PullLastUserSongRequest(channelID *string, userID *string) {
	songRequestInfo := GetSongRequest(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return
	}
	id := ""
	var date time.Time
	for _, request := range songRequestInfo.Requests {
		if request.UserID == *userID && request.Date.After(date) {
			id = request.VideoID
			date = request.Date
		}
	}
	if id != "" {
		PullUserSongRequest(channelID, &id, userID)
	}
}

func GetYoutubeVideoInfo(id *string) (*models.YoutubeVideo, error) {
	if Config.YoutubeKey == "" {
		return nil, errors.New("YT key is not set")
	}
	resp, error := httpclient.YoutubeVideo(id, &Config.YoutubeKey)
	if error != nil {
		return nil, error
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var ytVideo models.YoutubeVideo
	marshallError := json.NewDecoder(resp.Body).Decode(&ytVideo)
	if marshallError != nil {
		return nil, marshallError
	}
	return &ytVideo, nil
}
