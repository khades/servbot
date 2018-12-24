package channelBans

import (
	"github.com/globalsign/mgo"
	"time"

	"github.com/globalsign/mgo/bson"

)

type Service struct {
	// Own Fields
	collection *mgo.Collection
}
// Log logs user channel ban, limitting to last 500
func (service *Service) Log(userID *string, user *string, channelID *string, duration *int) {
	service.collection.Upsert(
		bson.M{"channelid": *channelID},
		bson.M{
			"$push": bson.M{
				"bans": bson.M{
					"$each": []ChannelBanRecord{ChannelBanRecord{
						User:          *user,
						UserID:        *userID,
						Date:          time.Now(),
                        BanLength: *duration}},
					"$sort":  bson.M{"date": -1},
					"$slice": 500}}})
}

// Search returns list of all bans on channel (limited to last 500)
func (service *Service) Get(channelID *string) (*ChannelBans, error) {
	var result = ChannelBans{}
	error := service.collection.Find(bson.M{"channelid": *channelID}).One(&result)
	return &result, error
}