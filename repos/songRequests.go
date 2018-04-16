package repos

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/khades/servbot/l10n"
	//"time"
	"encoding/json"
	"errors"
	"net/http"
	//"net/http/httputil"

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
	if strings.Contains(input, "youtube.com/watch?") {
		result := ""
		params := strings.Split(input, "youtube.com/watch?")[1]
		paramsSplit := strings.Split(params, "&")
		for _, param := range paramsSplit {
			paramSplit := strings.Split(param, "=")
			if paramSplit[0] == "v" {
				result = short(paramSplit[1], 11)
				break
			}
		}
		return result
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
	songRequestInfo := models.ChannelSongRequest{Settings: models.ChannelSongRequestSettings{PlaylistLength: 30, MaxVideoLength: 300, MaxRequestsPerUser: 2, MoreLikes: true}}
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
func GetTopRequest(channelID *string, lang string) models.CurrentSong {
	songRequestInfo := GetSongRequest(channelID)

	for _, request := range songRequestInfo.Requests {
		if request.Order == 1 {
			return models.CurrentSong{
				IsPlaying: true,
				Title:     request.Title,
				User:      request.User,
				Link:      "https://youtu.be/" + request.VideoID,
				Duration:  l10n.HumanizeDuration(request.Length, lang),
				Volume:    songRequestInfo.Settings.Volume,
				Count:     len(songRequestInfo.Requests)}
			break
		}
	}
	return models.CurrentSong{
		IsPlaying: false}
}

func SetSongRequestVolume(channelID *string, volume int) {
	db.C(songRequestCollectionName).Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$set": bson.M{"settings.volume": volume}})
	eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "volume:"+strconv.Itoa(volume))
}

func SetSongRequestVolumeNoEvent(channelID *string, volume int) {
	db.C(songRequestCollectionName).Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$set": bson.M{"settings.volume": volume}})
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
func AddSongRequest(user *string, userIsSub bool, userID *string, channelID *string, videoID *string) models.SongRequestAddResult {
	// logger := logrus.WithFields(logrus.Fields{
	// 	"package": "repos",
	// 	"feature": "songrequests",
	// 	"action":  "AddSongRequest"})
	songRequestInfo := GetSongRequest(channelID)
	channelInfo, channelInfoError := GetChannelInfo(channelID)

	if channelInfoError != nil || (songRequestInfo.Settings.AllowOffline == false && channelInfo.StreamStatus.Online == false) {
		return models.SongRequestAddResult{Offline: true}
	}
	if len(songRequestInfo.Requests) >= songRequestInfo.Settings.PlaylistLength {
		return models.SongRequestAddResult{PlaylistIsFull: true}
	}
	parsedVideoID := parseYoutubeLink(*videoID)

	for _, request := range songRequestInfo.Requests {
		if request.VideoID == parsedVideoID {
			return models.SongRequestAddResult{AlreadyInPlaylist: true, Title: request.Title, Length: request.Length}

		}
	}

	userRequestsCount := 0
	for _, request := range songRequestInfo.Requests {
		if request.UserID == *userID {
			userRequestsCount = userRequestsCount + 1
		}
	}

	if userRequestsCount >= songRequestInfo.Settings.MaxRequestsPerUser {
		return models.SongRequestAddResult{TooManyRequests: true}
	}

	if parsedVideoID == "" {
		return models.SongRequestAddResult{InvalidLink: true}
	}
	var songRequest models.SongRequest
	libraryItem, libraryError := getVideo(&parsedVideoID)
	if libraryError == nil {
		yTrestricted := false
		twitchRestricted := false
		channelRestricted := false
		tagRestricted := false
		bannedTag := ""
		for _, tag := range libraryItem.Tags {
			if tag.Tag == "youtuberestricted" {
				yTrestricted = true
				break
			}
			if tag.Tag == "twitchrestricted" {
				twitchRestricted = true
				break
			}
			if tag.Tag == *channelID+"-restricted" {
				channelRestricted = true
				break
			}
			for _, channelTag := range songRequestInfo.Settings.BannedTags {
				if tag.Tag == channelTag && tag.Tag != *channelID+"-restricted" {
					bannedTag = channelTag
					tagRestricted = true
					break
				}
			}
		}
		if yTrestricted == true {
			return models.SongRequestAddResult{YoutubeRestricted: true, Title: libraryItem.Title}
		}
		if tagRestricted == true {
			return models.SongRequestAddResult{TagRestricted: true, Title: libraryItem.Title, Tag: bannedTag}
		}
		if twitchRestricted == true {
			return models.SongRequestAddResult{TwitchRestricted: true, Title: libraryItem.Title}
		}
		if channelRestricted == true {
			return models.SongRequestAddResult{ChannelRestricted: true, Title: libraryItem.Title}
		}
	}
	if libraryError != nil || time.Now().Sub(libraryItem.LastCheck) > 3*60*time.Minute {
		video, videoError := getYoutubeVideoInfo(&parsedVideoID)
		if videoError != nil {
			return models.SongRequestAddResult{InternalError: true}
		}
		if len(video.Items) == 0 {
			return models.SongRequestAddResult{NothingFound: true}
		}
		duration, durationError := video.Items[0].ContentDetails.GetDuration()
		if durationError != nil {
			return models.SongRequestAddResult{InternalError: true}
		}
		likes, likesError := strconv.ParseInt(video.Items[0].Statistics.Likes, 10, 64)
		if likesError != nil {
			likes = 0
		}
		dislikes, dislikesError := strconv.ParseInt(video.Items[0].Statistics.Dislikes, 10, 64)
		if dislikesError != nil {
			dislikes = 0
		}
		addVideoToLibrary(&parsedVideoID, &video.Items[0].Snippet.Title, duration, video.Items[0].Statistics.GetViewCount(), likes, dislikes)
		songRequest = models.SongRequest{
			User:     *user,
			UserID:   *userID,
			Date:     time.Now(),
			VideoID:  parsedVideoID,
			Length:   *duration,
			Order:    len(songRequestInfo.Requests) + 1,
			Title:    video.Items[0].Snippet.Title,
			Likes:    likes,
			Dislikes: dislikes,
			Views:    video.Items[0].Statistics.GetViewCount()}
	} else {

		songRequest = models.SongRequest{
			User:     *user,
			UserID:   *userID,
			Date:     time.Now(),
			VideoID:  parsedVideoID,
			Length:   libraryItem.Length,
			Order:    len(songRequestInfo.Requests) + 1,
			Title:    libraryItem.Title,
			Likes:    libraryItem.Likes,
			Dislikes: libraryItem.Dislikes,
			Views:    libraryItem.Views}
	}

	if songRequest.Length.Seconds() > float64(songRequestInfo.Settings.MaxVideoLength) {
		return models.SongRequestAddResult{TooLong: true, Title: songRequest.Title, Length: songRequest.Length}

	}
	if songRequest.Views < songRequestInfo.Settings.VideoViewLimit {
		return models.SongRequestAddResult{TooLittleViews: true, Title: songRequest.Title, Length: songRequest.Length}

	}
	if songRequestInfo.Settings.MoreLikes == true && songRequest.Dislikes > songRequest.Likes {
		return models.SongRequestAddResult{MoreDislikes: true, Title: songRequest.Title, Length: songRequest.Length}
	}

	PushSongRequest(channelID, &songRequest)

	eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")

	return models.SongRequestAddResult{Success: true, Title: songRequest.Title, Length: songRequest.Length}
}

// PullSongRequest removes songrequest, specified by youtube video ID on specified channel
func PullSongRequest(channelID *string, videoID *string) {
	songRequestInfo := GetSongRequest(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return
	}
	newRequests, pulledItem := songRequestInfo.Requests.PullOneRequest(videoID)

	if pulledItem != nil {
		putRequests(channelID, newRequests)
		eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
	}
}

// func SetSongRequestRestricted(channelID *string, videoID *string) {
// 	AddTagToVideo(videoID, "youtuberestricted",)
// 	PullSongRequest(channelID, videoID)
// }
// // PullUserSongRequest removes specified user request, specified by youtube video ID on specified channel
// func PullUserSongRequest(channelID *string, videoID *string, userID *string) {
// 	db.C(songRequestCollectionName).Update(
// 		bson.M{
// 			"channelid": *channelID}, bson.M{"$pull": bson.M{"requests": bson.M{"userid:": *userID, "videoid": *videoID}}})
// 	eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
// }

// PullLastUserSongRequest removes last specified user request on specified channel
func PullLastUserSongRequest(channelID *string, userID *string) (*models.SongRequest, bool) {
	songRequestInfo := GetSongRequest(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return nil, false
	}
	newRequests, pulledItem := songRequestInfo.Requests.PullUsersLastRequest(userID)

	if pulledItem != nil {

		putRequests(channelID, newRequests)
		eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
		return pulledItem, true
	}
	return nil, false

}

func PushSettings(channelID *string, settings *models.ChannelSongRequestSettings) {
	db.C(songRequestCollectionName).Update(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"settings": *settings}})
}

// BubbleUpVideo sets order of song to 1, and increases order of other songs
func BubbleUpVideo(channelID *string, videoID *string) bool {
	songRequestInfo := GetSongRequest(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return false
	}
	newRequests, changed := songRequestInfo.Requests.BubbleVideoUp(videoID, 1)
	if changed == true {
		putRequests(channelID, newRequests)
		eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
	}
	return changed
}

// BubbleUpVideoToSecond sets order of song to 2
func BubbleUpVideoToSecond(channelID *string, videoID *string) bool {
	songRequestInfo := GetSongRequest(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return false
	}
	newRequests, changed := songRequestInfo.Requests.BubbleVideoUp(videoID, 2)
	if changed == true {
		putRequests(channelID, newRequests)
		eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "update")
	}
	return changed
}

func putRequests(channelID *string, requests models.SongRequests) {
	db.C(songRequestCollectionName).Update(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"requests": requests}})
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
