package songRequest

import (
	"strconv"

	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/videoLibrary"
	"github.com/khades/servbot/youtubeAPI"

	"github.com/khades/servbot/l10n"
	"github.com/sirupsen/logrus"

	//"time"

	//"net/http/httputil"

	"time"

	"github.com/BurntSushi/locker"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/eventbus"
)

type Service struct {
	collection          *mgo.Collection
	youtubeAPIClient    *youtubeAPI.Client
	channelInfoService  *channelInfo.Service
	videoLibraryService *videoLibrary.Service
	eventBus            EventBus.Bus
}

// Get gets full songrequest info for specified channel
func (service *Service) Get(channelID *string) *ChannelSongRequest {
	songRequestInfo := ChannelSongRequest{Settings: ChannelSongRequestSettings{PlaylistLength: 30, MaxVideoLength: 300, MaxRequestsPerUser: 2, MoreLikes: true}}
	service.collection.Find(
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

func (service *Service) formCurrentSong(channelID *string, volume int, songRequests SongRequests) (*channelInfo.CurrentSong, error) {
	for _, request := range songRequests {
		if request.Order == 1 {
			channelInfoStruct, err := service.channelInfoService.Get(channelID)
			if err != nil {
				return nil, err
			}
			return &channelInfo.CurrentSong{
				IsPlaying: true,
				Title:     request.Title,
				User:      request.User,
				Link:      "https://youtu.be/" + request.VideoID,
				Duration:  l10n.HumanizeDurationFull(request.Length, channelInfoStruct.Lang, true),
				Volume:    volume,
				Count:     len(songRequests),
				ID:        request.VideoID}, nil
			break
		}
	}
	return &channelInfo.CurrentSong{
		IsPlaying: false}, nil
}

func (service *Service) GetLast(channelID *string, lang string) channelInfo.CurrentSong {
	songRequestInfo := service.Get(channelID)

	for _, request := range songRequestInfo.Requests {
		if request.Order == 1 {
			return channelInfo.CurrentSong{
				IsPlaying: true,
				Title:     request.Title,
				User:      request.User,
				Link:      "https://youtu.be/" + request.VideoID,
				Duration:  l10n.HumanizeDurationFull(request.Length, lang, true),
				Volume:    songRequestInfo.Settings.Volume,
				Count:     len(songRequestInfo.Requests),
				ID:        request.VideoID}
			break
		}
	}
	return channelInfo.CurrentSong{
		IsPlaying: false}
}

func (service *Service) SetVolume(channelID *string, volume int) {
	service.collection.Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$set": bson.M{"settings.volume": volume}})
	service.eventBus.Publish(eventbus.Songrequest(channelID), "volume:"+strconv.Itoa(volume))
}

func (service *Service) SetVolumeNoEvent(channelID *string, volume int) {
	service.collection.Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$set": bson.M{"settings.volume": volume}})
}

func (service *Service) push(channelID *string, request *SongRequest) {
	service.collection.Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$push": bson.M{"requests": *request}})
}

// SetSettings updates songrequest settings for specified channel
func (service *Service) SetSettings(channelID *string, settings *ChannelSongRequestSettings) {
	service.collection.Upsert(
		bson.M{
			"channelid": *channelID}, bson.M{"$set": bson.M{"settings": *settings}})

}

// Add processes youtube video link before pushing it to songrequest database
func (service *Service) Add(user *string, userIsSub bool, userID *string, channelID *string, videoID *string) SongRequestAddResult {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "songrequests",
		"action":  "Add"})
	locker.Lock("sr" + *channelID)
	logger.Debugf("Setting Songrequest lock for channel %s", *channelID)

	defer locker.Unlock("sr" + *channelID)
	defer logger.Debugf("Releasing Songrequest lock for channel %s", *channelID)

	songRequestInfo := service.Get(channelID)
	channelInfo, channelInfoError := service.channelInfoService.Get(channelID)

	if channelInfoError != nil || (songRequestInfo.Settings.AllowOffline == false && channelInfo.StreamStatus.Online == false) {
		return SongRequestAddResult{Offline: true}
	}

	if len(songRequestInfo.Requests) >= songRequestInfo.Settings.PlaylistLength {
		return SongRequestAddResult{PlaylistIsFull: true}
	}

	parsedVideoID, parsedVideoIsID := parseYoutubeLink(*videoID)

	if parsedVideoIsID == true {
		for _, request := range songRequestInfo.Requests {
			if request.VideoID == parsedVideoID {
				return SongRequestAddResult{AlreadyInPlaylist: true, Title: request.Title, Length: request.Length}
			}
		}
	}

	userRequestsCount := 0
	for _, request := range songRequestInfo.Requests {
		if request.UserID == *userID {
			userRequestsCount = userRequestsCount + 1
		}
	}

	if userRequestsCount >= songRequestInfo.Settings.MaxRequestsPerUser {
		return SongRequestAddResult{TooManyRequests: true}
	}

	var songRequest SongRequest
	var libraryItem = &videoLibrary.SongRequestLibraryItem{}
	var libraryError error

	if parsedVideoIsID == true {
		libraryItem, libraryError = service.videoLibraryService.GetVideo(&parsedVideoID)
	}

	if parsedVideoIsID == false || libraryError == nil {
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
			return SongRequestAddResult{YoutubeRestricted: true, Title: libraryItem.Title}
		}
		if tagRestricted == true {
			return SongRequestAddResult{TagRestricted: true, Title: libraryItem.Title, Tag: bannedTag}
		}
		if twitchRestricted == true {
			return SongRequestAddResult{TwitchRestricted: true, Title: libraryItem.Title}
		}
		if channelRestricted == true {
			return SongRequestAddResult{ChannelRestricted: true, Title: libraryItem.Title}
		}
	}
	if libraryError != nil || time.Now().Sub(libraryItem.LastCheck) > 3*60*time.Minute {
		var videoError error
		var video = &youtubeAPI.YoutubeVideo{}
		if parsedVideoIsID == true {
			video, videoError = service.youtubeAPIClient.Get(&parsedVideoID)
		}
		// if videoError != nil {
		// 	logger.Infof("Youtube error: %s", videoError.Error())
		// 	return models.SongRequestAddResult{InternalError: true}
		// }
		if parsedVideoIsID == false || videoError != nil || len(video.Items) == 0 {
			var videoStringError error
			video, videoStringError = service.youtubeAPIClient.Search(&parsedVideoID)
			if videoStringError != nil {
				logger.Infof("Youtube error: %s", videoError.Error())
				return SongRequestAddResult{InternalError: true}
			}
			if len(video.Items) == 0 {
				return SongRequestAddResult{NothingFound: true}
			}
		}

		if parsedVideoIsID == false {
			for _, request := range songRequestInfo.Requests {
				if request.VideoID == video.Items[0].ID {
					return SongRequestAddResult{AlreadyInPlaylist: true, Title: request.Title, Length: request.Length}
				}
			}
		}

		duration, durationError := video.Items[0].ContentDetails.GetDuration()
		if durationError != nil {
			logger.Infof("Youtube error: %s", videoError.Error())
			return SongRequestAddResult{InternalError: true}
		}
		likes, likesError := strconv.ParseInt(video.Items[0].Statistics.Likes, 10, 64)
		if likesError != nil {
			likes = 0
		}
		dislikes, dislikesError := strconv.ParseInt(video.Items[0].Statistics.Dislikes, 10, 64)
		if dislikesError != nil {
			dislikes = 0
		}
		service.videoLibraryService.Add(&video.Items[0].ID, &video.Items[0].Snippet.Title, duration, video.Items[0].Statistics.GetViewCount(), likes, dislikes)
		songRequest = SongRequest{
			User:     *user,
			UserID:   *userID,
			Date:     time.Now(),
			VideoID:  video.Items[0].ID,
			Length:   *duration,
			Order:    len(songRequestInfo.Requests) + 1,
			Title:    video.Items[0].Snippet.Title,
			Likes:    likes,
			Dislikes: dislikes,
			Views:    video.Items[0].Statistics.GetViewCount()}
	} else {

		songRequest = SongRequest{
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
		return SongRequestAddResult{TooLong: true, Title: songRequest.Title, Length: songRequest.Length}

	}
	if songRequest.Length.Seconds() < float64(songRequestInfo.Settings.MinVideoLength) {
		return SongRequestAddResult{TooShort: true, Title: songRequest.Title, Length: songRequest.Length}

	}
	if songRequest.Views < songRequestInfo.Settings.VideoViewLimit {
		return SongRequestAddResult{TooLittleViews: true, Title: songRequest.Title, Length: songRequest.Length}

	}
	if songRequestInfo.Settings.MoreLikes == true && songRequest.Dislikes > songRequest.Likes {
		return SongRequestAddResult{MoreDislikes: true, Title: songRequest.Title, Length: songRequest.Length}
	}

	service.push(channelID, &songRequest)

	service.eventBus.Publish(eventbus.Songrequest(channelID), "update")

	return SongRequestAddResult{Success: true, Title: songRequest.Title, Length: songRequest.Length}
}

// Pull removes songrequest, specified by youtube video ID on specified channel
func (service *Service) Pull(channelID *string, videoID *string) {
	locker.Lock("sr" + *channelID)
	defer locker.Unlock("sr" + *channelID)
	songRequestInfo := service.Get(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return
	}
	newRequests, pulledItem := songRequestInfo.Requests.PullOneRequest(videoID)

	if pulledItem != nil {
		service.put(channelID, newRequests, songRequestInfo.Settings.Volume)
		service.eventBus.Publish(eventbus.Songrequest(channelID), "update")
	}
}

// func SetSongRequestRestricted(channelID *string, videoID *string) {
// 	PushTag(videoID, "youtuberestricted",)
// 	Pull(channelID, videoID)
// }
// // PullUserSongRequest removes specified user request, specified by youtube video ID on specified channel
// func PullUserSongRequest(channelID *string, videoID *string, userID *string) {
// 	service.collection.Update(
// 		bson.M{
// 			"channelid": *channelID}, bson.M{"$pull": bson.M{"requests": bson.M{"userid:": *userID, "videoid": *videoID}}})
// 	service.eventBus.Publish(eventbus.Songrequest(channelID), "update")
// }

// PullLastUser removes last specified user request on specified channel
func (service *Service) PullLastUser(channelID *string, userID *string) (*SongRequest, bool) {
	locker.Lock("sr" + *channelID)
	defer locker.Unlock("sr" + *channelID)
	songRequestInfo := service.Get(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return nil, false
	}
	newRequests, pulledItem := songRequestInfo.Requests.PullUsersLastRequest(userID)

	if pulledItem != nil {

		service.put(channelID, newRequests, songRequestInfo.Settings.Volume)
		service.eventBus.Publish(eventbus.Songrequest(channelID), "update")
		return pulledItem, true
	}
	return nil, false

}

//func (service *Service) PushSettings(channelID *string, settings *ChannelSongRequestSettings) {
//	service.collection.Update(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"settings": *settings}})
//}

// BubbleUp sets order of song to 1, and increases order of other songs
func (service *Service) BubbleUp(channelID *string, videoID *string) bool {
	locker.Lock("sr" + *channelID)
	defer locker.Unlock("sr" + *channelID)
	songRequestInfo := service.Get(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return false
	}
	newRequests, changed := songRequestInfo.Requests.BubbleVideoUp(videoID, 1)
	if changed == true {
		service.put(channelID, newRequests, songRequestInfo.Settings.Volume)
		service.eventBus.Publish(eventbus.Songrequest(channelID), "update")
	}
	return changed
}

// BubbleUpToSecond sets order of song to 2
func (service *Service) BubbleUpToSecond(channelID *string, videoID *string) bool {
	locker.Lock("sr" + *channelID)
	defer locker.Unlock("sr" + *channelID)
	songRequestInfo := service.Get(channelID)
	if len(songRequestInfo.Requests) == 0 {
		return false
	}
	newRequests, changed := songRequestInfo.Requests.BubbleVideoUp(videoID, 2)
	if changed == true {
		service.put(channelID, newRequests, songRequestInfo.Settings.Volume)
		service.eventBus.Publish(eventbus.Songrequest(channelID), "update")
	}
	return changed
}

func (service *Service) put(channelID *string, requests SongRequests, volume int) {
	currentSong, _ := service.formCurrentSong(channelID,  volume, requests)
	service.channelInfoService.PutCurrentSong(channelID, currentSong)
	service.collection.Update(bson.M{"channelid": *channelID}, bson.M{"$set": bson.M{"requests": requests}})
}

func (service *Service) PushTag(videoID *string, tag string, userID string, user string) []TaggedVideoResult {
	var channels []ChannelSongRequest
	var tagResults []TaggedVideoResult
	track := service.videoLibraryService.PushTag(videoID, tag, userID, user)
	error := service.collection.Find(bson.M{"requests.videoid": *videoID}).All(&channels)
	if error != nil {
		return []TaggedVideoResult{}
	}
	for _, channel := range channels {
		pull := false
		if tag == "youtuberestricted" {
			pull = true
			tagResults = append(tagResults, TaggedVideoResult{
				Title:                    track.Title,
				ChannelID:                channel.ChannelID,
				RemovedYoutubeRestricted: true})
		}
		if tag == "twitchrestricted" {
			pull = true
			tagResults = append(tagResults, TaggedVideoResult{
				Title:                   track.Title,
				ChannelID:               channel.ChannelID,
				RemovedTwitchRestricted: true})
		}
		if tag == channel.ChannelID+"-restricted" {
			pull = true
			tagResults = append(tagResults, TaggedVideoResult{
				Title:                    track.Title,
				ChannelID:                channel.ChannelID,
				RemovedChannelRestricted: true})

			service.videoLibraryService.IncrementBannedTracks(&channel.ChannelID)

		}
		if channel.Settings.SkipIfTagged == true {
			for _, channelTag := range channel.Settings.BannedTags {
				if tag == channelTag && tag != channel.ChannelID+"-restricted" {
					tagResults = append(tagResults, TaggedVideoResult{
						Title:                track.Title,
						ChannelID:            channel.ChannelID,
						RemovedTagRestricted: true,
						Tag:                  tag})
					pull = true
					break
				}
			}
		}
		if pull == false {
			service.collection.Update(bson.M{"channelid": channel.ChannelID, "requests": bson.M{"videoid": *videoID, "tags.tag": bson.M{"$ne": tag}}}, bson.M{"$push": bson.M{"requests.$.tags": TagRecord{User: user, UserID: userID, Tag: tag}}})
		} else {
			newRequests, _ := channel.Requests.PullOneRequest(videoID)
			service.put(&channel.ChannelID, newRequests, channel.Settings.Volume)
		}
		service.eventBus.Publish(eventbus.Songrequest(&channel.ChannelID), "update")

	}

	return tagResults
}
