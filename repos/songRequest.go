package repos

import (
	//"time"
	"sort"
	"errors"
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var songRequestCollectionName = "songrequests"

func AddSongRequest(user *string, userID *string, channelID *string, videoID *string) error {
	songRequestInfo := models.ChannelSongRequest {}
	error := Db.C(songRequestCollectionName).Find(
		bson.M{
			"channelid": bson.ObjectId(*channelID)}).One(&songRequestInfo)
	if error != nil {
		return error
	}
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
	if userRequestsCount >=  songRequestInfo.Settings.MaxRequestsPerUser {
		return errors.New("Too many requests per user")
	}
	return nil
}

func RemoveSongRequest(channelID *string, videoID *string) {
	Db.C(songRequestCollectionName).Update(
		bson.M{
			"channelid": bson.ObjectId(*channelID)},
		bson.M{"$set": bson.M{
			"inqueue":    false,
			"playingnow": false,
		}})

}

// func SetTrackPlayed(id string) {

// 	Db.C("songRequests").Update(
// 		bson.M{
// 			"_id": bson.ObjectId(id)},
// 		bson.M{"$set": bson.M{
// 			"inqueue":    false,
// 			"playingnow": false,
// 		}})

// }
