package repos

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
)

var videolibraryCollection = "videolibrary"
var videolibraryCount = -1
var videolibraryBannedCountPerChannel = make(map[string]int)

func GetCount() int {
	if videolibraryCount == -1 {
		result, _ := db.C(videolibraryCollection).Count()
		videolibraryCount = result
		return result
	}
	return videolibraryCount
}
func getBannedTracksCountForChannel(channelID *string) int {
	bannedVideos, bannedVideosFound := videolibraryBannedCountPerChannel[*channelID]
	if bannedVideosFound == false {
		count, _ := db.C(videolibraryCollection).Find(bson.M{"tags.tag": *channelID + "-restricted"}).Count()
		videolibraryBannedCountPerChannel[*channelID] = count
		return count
	}
	return bannedVideos

}

// func getVideoItem(videoID *string) *models.SongRequestLibraryResponse {
// 	logger := logrus.WithFields(logrus.Fields{
// 		"package": "repos",
// 		"feature": "songrequests",
// 		"action":  "getVideoItem"})
// 	result := &models.SongRequestLibraryResponse{}
// 	parsedVideoID, parsedVideoIsID := parseYoutubeLink(*videoID)
// 	var libraryItem = &models.SongRequestLibraryItem{}
// 	var libraryError error

// 	if parsedVideoIsID == true {
// 		libraryItem, libraryError = getVideo(&parsedVideoID)
// 	}
// 	if libraryError == nil && time.Now().Sub(libraryItem.LastCheck) < 3*60*time.Minute {
// 		result.VideoID = libraryItem.VideoID
// 		result.Item = libraryItem
// 		return result
// 	}

// 	var videoError error
// 	var video = &models.YoutubeVideo{}
// 	if parsedVideoIsID == true {
// 		result.VideoID = parsedVideoID
// 		video, videoError = getYoutubeVideoInfo(&parsedVideoID)
// 	}

// 	if parsedVideoIsID == false || videoError != nil || len(video.Items) == 0 {
// 		var videoStringError error
// 		video, videoStringError = getYoutubeVideoInfoByString(&parsedVideoID)
// 		if videoStringError != nil {
// 			logger.Infof("Youtube error: %s", videoError.Error())
// 			result.InternalError = true
// 			return result
// 		}
// 		if len(video.Items) == 0 {
// 			result.VideoDoesntExist = true
// 			return result
// 		}
// 	}

// 	duration, durationError := video.Items[0].ContentDetails.GetDuration()
// 	if durationError != nil {
// 		logger.Infof("Youtube error: %s", videoError.Error())
// 		result.InternalError = true
// 		return result
// 	}
// 	likes, likesError := strconv.ParseInt(video.Items[0].Statistics.Likes, 10, 64)
// 	if likesError != nil {
// 		likes = 0
// 	}
// 	dislikes, dislikesError := strconv.ParseInt(video.Items[0].Statistics.Dislikes, 10, 64)
// 	if dislikesError != nil {
// 		dislikes = 0
// 	}
// 	addVideoToLibrary(&video.Items[0].ID, &video.Items[0].Snippet.Title, duration, video.Items[0].Statistics.GetViewCount(), likes, dislikes)
// 	libraryItem.Dislikes = dislikes
// 	libraryItem.LastCheck = time.Now()
// 	libraryItem.Likes = likes
// 	libraryItem.Length = *duration
// 	libraryItem.Title = video.Items[0].Snippet.Title
// 	libraryItem.Views = video.Items[0].Statistics.GetViewCount()
// 	libraryItem.VideoID = video.Items[0].ID
// 	result.VideoID = video.Items[0].ID

// 	result.Item = libraryItem
// 	return result
// }

func getVideo(videoID *string) (*models.SongRequestLibraryItem, error) {
	var result models.SongRequestLibraryItem
	err := db.C(videolibraryCollection).Find(bson.M{"videoid": *videoID}).One(&result)
	return &result, err
}

func addVideoToLibrary(videoID *string, title *string, duration *time.Duration, views int64, likes int64, dislikes int64) {

	changeInfo, err := db.C(videolibraryCollection).Upsert(bson.M{"videoid": *videoID}, bson.M{"$set": bson.M{
		"length":    *duration,
		"lastcheck": time.Now(),
		"likes":     likes,
		"views":     views,
		"dislikes":  dislikes,
		"title":     *title}})
	if err == nil && changeInfo.Matched == 0 && videolibraryCount != -1 {
		videolibraryCount = videolibraryCount + 1
	}
}

func GetBannedTracksForChannel(channelID *string, page int) ([]models.SongRequestLibraryItem, int, error) {
	pageSize := 25
	result := []models.SongRequestLibraryItem{}
	error := db.C(videolibraryCollection).Find(bson.M{"tags.tag": *channelID + "-restricted"}).Sort("-_id").Skip((page - 1) * pageSize).Limit(pageSize).All(&result)
	if error != nil {
		return result, 0, error
	}
	return result, getBannedTracksCountForChannel(channelID), error
}

func GetVideoLibraryItems(page int) ([]models.SongRequestLibraryItem, int, error) {
	pageSize := 25
	result := []models.SongRequestLibraryItem{}
	error := db.C(videolibraryCollection).Find(nil).Sort("-_id").Skip((page - 1) * pageSize).Limit(pageSize).All(&result)
	if error != nil {
		return result, 0, error
	}
	return result, GetCount(), error
}

func PullTagFromVideo(videoID *string, tag string) {
	db.C(videolibraryCollection).Update(bson.M{"videoid": *videoID}, bson.M{
		"$pull": bson.M{
			"tags": bson.M{"tag": tag}}})
}

func AddTagToVideo(videoID *string, tag string, userID string, user string) []models.TaggedVideoResult {
	var channels []models.ChannelSongRequest
	var tagResults []models.TaggedVideoResult
	var track models.SongRequestLibraryItem
	db.C(videolibraryCollection).Find(bson.M{"videoid": *videoID, "tags.tag": bson.M{"$ne": tag}}).Apply(mgo.Change{
		Update: bson.M{"$push": bson.M{
			"tags": models.TagRecord{User: user, UserID: userID, Tag: tag}}}}, &track)
	error := db.C(songRequestCollectionName).Find(bson.M{"requests.videoid": *videoID}).All(&channels)
	if error != nil {
		return []models.TaggedVideoResult{}
	}
	for _, channel := range channels {
		pull := false
		if tag == "youtuberestricted" {
			pull = true
			tagResults = append(tagResults, models.TaggedVideoResult{
				Title:                    track.Title,
				ChannelID:                channel.ChannelID,
				RemovedYoutubeRestricted: true})
		}
		if tag == "twitchrestricted" {
			pull = true
			tagResults = append(tagResults, models.TaggedVideoResult{
				Title:                   track.Title,
				ChannelID:               channel.ChannelID,
				RemovedTwitchRestricted: true})
		}
		if tag == channel.ChannelID+"-restricted" {
			pull = true
			tagResults = append(tagResults, models.TaggedVideoResult{
				Title:                    track.Title,
				ChannelID:                channel.ChannelID,
				RemovedChannelRestricted: true})
			bannedVideos, bannedVideosFound := videolibraryBannedCountPerChannel[channel.ChannelID]
			if bannedVideosFound == false {
				videolibraryBannedCountPerChannel[channel.ChannelID] = 1
			} else {
				videolibraryBannedCountPerChannel[channel.ChannelID] = bannedVideos + 1
			}
		}
		if channel.Settings.SkipIfTagged == true {
			for _, channelTag := range channel.Settings.BannedTags {
				if tag == channelTag && tag != channel.ChannelID+"-restricted" {
					tagResults = append(tagResults, models.TaggedVideoResult{
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
			db.C(songRequestCollectionName).Update(bson.M{"channelid": channel.ChannelID, "requests": bson.M{"videoid": *videoID, "tags.tag": bson.M{"$ne": tag}}}, bson.M{"$push": bson.M{"requests.$.tags": models.TagRecord{User: user, UserID: userID, Tag: tag}}})
		} else {
			newRequests, _ := channel.Requests.PullOneRequest(videoID)
			putRequests(&channel.ChannelID, newRequests)
		}
		eventbus.EventBus.Publish(eventbus.Songrequest(&channel.ChannelID), "update")

	}

	return tagResults
}
