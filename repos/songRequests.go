package repos

import (
	"strings"
	"unicode/utf8"
	//"time"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
)

var songRequestCollectionName = "songrequests"

func short(s string, i int) string {
	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}

func parseYoutubeLink(input string) string {

	if utf8.RuneCountInString(input) == 11 {
		return input

	}
	if strings.Contains(input, "youtube.com/watch?v=") {
		return short(strings.Split(input, "youtube.com/watch?v=")[1], 11)
	}
	if strings.Contains(input, "youtube.com/v/") {
		return short(strings.Split(input, "youtube.com/v/")[1], 11)

	}
	if strings.Contains(input, "youtu.be/") {
		return short(strings.Split(input, "youtu.be/")[1], 11)
	}
	return ""
}

// GetSongRequest gets full songrequest info for specified channel
func GetSongRequest(channelID *string) *models.ChannelSongRequest {
	songRequestInfo := models.ChannelSongRequest{Settings: models.ChannelSongRequestSettings{PlaylistLength: 30, MaxVideoLength: 300, MaxRequestsPerUser: 2}}
	db.C(songRequestCollectionName).Find(
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
		songRequestInfo.Settings.VideoViewLimit = 2000

	}
	return &songRequestInfo
}

// PushSongRequest pushes songrequest for specified channel
func PushSongRequest(channelID *string, request *models.SongRequest) {
	db.C(songRequestCollectionName).Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$push": bson.M{"requests": *request}})
}

// PushSongRequestSettings updates songrequest settings for specified channel
func PushSongRequestSettings(channelID *string, settings *models.ChannelSongRequestSettings) {
	db.C(songRequestCollectionName).Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$set": bson.M{"settings": *settings}})
}

// AddSongRequest processes youtube video link before pushing it to songrequest database
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

	parsedVideoID := parseYoutubeLink(*videoID)

	if (parsedVideoID == "") {
		return errors.New("Invalid link")
	}

	video, videoError := getYoutubeVideoInfo(&parsedVideoID)
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
		return errors.New("Too little views on video, got " + string(video.Items[0].Statistics.ViewCount))

	}
	songRequest := models.SongRequest{
		User:    *user,
		UserID:  *userID,
		Date:    time.Now(),
		VideoID: parsedVideoID,
		Length:  *duration,
		Title:   video.Items[0].Snippet.Title}

	PushSongRequest(channelID, &songRequest)

	eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")

	return errors.New("Video " + video.Items[0].Snippet.Title + " , length: " + duration.String())
}

// PullSongRequest removes songrequest, specified by youtube video ID on specified channel
func PullSongRequest(channelID *string, videoID *string) {
	db.C(songRequestCollectionName).Update(
		bson.M{
			"channelid": *channelID}, bson.M{"$pull": bson.M{"requests": bson.M{"videoid": *videoID}}})
	eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
}

// PullUserSongRequest removes specified user request, specified by youtube video ID on specified channel
func PullUserSongRequest(channelID *string, videoID *string, userID *string) {
	db.C(songRequestCollectionName).Update(
		bson.M{
			"channelid": *channelID}, bson.M{"$pull": bson.M{"requests": bson.M{"userid:": *userID, "videoid": *videoID}}})
	eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
}

// PullLastUserSongRequest removes last specified user request on specified channel
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
		eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
	}
}

func getYoutubeVideoInfo(id *string) (*models.YoutubeVideo, error) {
	if Config.YoutubeKey == "" {
		return nil, errors.New("YT key is not set")
	}
	resp, error := getYoutubeVideo(id)
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

func getYoutubeVideo(id *string) (*http.Response, error) {
	url := "https://content.googleapis.com/youtube/v3/videos?id=" + *id + "&part=snippet%2CcontentDetails%2Cstatistics&key=" + Config.YoutubeKey
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	return client.Get(url)
}
