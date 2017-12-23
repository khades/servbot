package repos
import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
)

var channelBansCollectionName = "channelBans"


func LogChannelBan(userID *string, user *string, channelID *string, duration *int) {
	Db.C(channelBansCollectionName).Upsert(
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

func GetChannelBans(channelID *string) (*models.ChannelBans, error) {
	var result = models.ChannelBans{}
	error := Db.C(channelBansCollectionName).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}