package repos
import (
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/khades/servbot/models"
)

var channelBansCollectionName = "channelBans"

// LogChannelBan logs user channel ban, limitting to last 500 
func LogChannelBan(userID *string, user *string, channelID *string, duration *int) {
	db.C(channelBansCollectionName).Upsert(
		bson.M{"channelid": *channelID},
		bson.M{
			"$push": bson.M{
				"bans": bson.M{
					"$each": []models.ChannelBanRecord{models.ChannelBanRecord{
						User:          *user,
						UserID:        *userID,
						Date:          time.Now(),
                        BanLength: *duration}},
					"$sort":  bson.M{"date": -1},
					"$slice": 500}}})
}

// GetChannelBans returns list of all bans on channel (limited to last 500)
func GetChannelBans(channelID *string) (*models.ChannelBans, error) {
	var result = models.ChannelBans{}
	error := db.C(channelBansCollectionName).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}