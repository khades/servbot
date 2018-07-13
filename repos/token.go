package repos

import (
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/utils"
)

const tokenCollection = "tokens"

func GetChannelToken(channelID string) (string, error) {
	result := models.Token{}
	err := db.C(tokenCollection).Find(bson.M{"channelid": channelID }).One(&result)
	if err.Error() == "not found" {
		return RandomizeChannelToken(channelID), nil
	}
	return result.Token, err
}

func RandomizeChannelToken(channelID string) string {
	token := utils.RandomString(24)
	db.C(tokenCollection).Upsert(bson.M{"channeliid": channelID }, bson.M{"$set": bson.M{"token": token}})
	return token
}
