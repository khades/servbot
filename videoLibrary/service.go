package videoLibrary

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	count                 int
	collection            *mgo.Collection
	bannedCountPerChannel map[string]int
}

func (service *Service) getCount() int {
	if service.count == -1 {
		result, _ := service.collection.Count()
		service.count = result
		return result
	}
	return service.count
}

func (service *Service) getBannedTrackCount(channelID *string) int {
	bannedVideos, bannedVideosFound := service.bannedCountPerChannel[*channelID]
	if !bannedVideosFound {
		count, _ := service.collection.Find(bson.M{"tags.tag": *channelID + "-restricted"}).Count()
		service.bannedCountPerChannel[*channelID] = count
		return count
	}
	return bannedVideos

}

func (service *Service) IncrementBannedTracks(channelID *string) {
	bannedTracks, bannedVideosFound := service.bannedCountPerChannel[*channelID]
	if !bannedVideosFound {
		service.bannedCountPerChannel[*channelID] = 1
	} else {
		service.bannedCountPerChannel[*channelID] = bannedTracks + 1
	}
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

func (service *Service) GetVideo(videoID *string) (*SongRequestLibraryItem, error) {
	var result SongRequestLibraryItem
	err := service.collection.Find(bson.M{"videoid": *videoID}).One(&result)
	return &result, err
}

func (service *Service) Add(videoID *string, title *string, duration *time.Duration, views int64, likes int64, dislikes int64) {

	changeInfo, err := service.collection.Upsert(bson.M{"videoid": *videoID}, bson.M{"$set": bson.M{
		"length":    *duration,
		"lastcheck": time.Now(),
		"likes":     likes,
		"views":     views,
		"dislikes":  dislikes,
		"title":     *title}})
	if err == nil && changeInfo.Matched == 0 && service.count != -1 {
		service.count = service.count + 1
	}
}

func (service *Service) ListBannedTracks(channelID *string, page int) ([]SongRequestLibraryItem, int, error) {
	pageSize := 25
	result := []SongRequestLibraryItem{}
	error := service.collection.Find(bson.M{"tags.tag": *channelID + "-restricted"}).Sort("-_id").Skip((page - 1) * pageSize).Limit(pageSize).All(&result)
	if error != nil {
		return result, 0, error
	}
	return result, service.getBannedTrackCount(channelID), error
}

func (service *Service) List(page int) ([]SongRequestLibraryItem, int, error) {
	pageSize := 25
	result := []SongRequestLibraryItem{}
	error := service.collection.Find(nil).Sort("-_id").Skip((page - 1) * pageSize).Limit(pageSize).All(&result)
	if error != nil {
		return result, 0, error
	}
	return result, service.getCount(), error
}

func (service *Service) PullTag(videoID *string, tag string) {
	service.collection.Update(bson.M{"videoid": *videoID}, bson.M{
		"$pull": bson.M{
			"tags": bson.M{"tag": tag}}})
}

func (service *Service) PushTag(videoID *string, tag string, userID string, user string) *SongRequestLibraryItem {
	var track SongRequestLibraryItem

	service.collection.Find(bson.M{"videoid": *videoID, "tags.tag": bson.M{"$ne": tag}}).Apply(mgo.Change{
		Update: bson.M{"$push": bson.M{
			"tags": TagRecord{User: user, UserID: userID, Tag: tag}}}}, &track)
	return &track
}
