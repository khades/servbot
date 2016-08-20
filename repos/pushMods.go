package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// PushMods updates list of mods on channel
func PushMods(channel string, mods []string) {
	Db.C("channelInfo").Upsert(models.ChannelSelector{Channel: channel}, bson.M{"$set": bson.M{"mods": mods}})
}
