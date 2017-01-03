package repos

import (
	"time"

	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// GetExpiredSongRequests checks for expired song requests to switch tracks then
func GetExpiredSongRequests() (*[]models.SongRequest, error) {
	result := []models.SongRequest{}

	error := Db.C("songRequests").Find(bson.M{
		"paused":     false,
		"playingnow": true,
		"estimatedendtime": bson.M{
			"$lt": time.Now()}}).All(&result)
	return &result, error
}
