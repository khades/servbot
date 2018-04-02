package repos

import (
	"time"

	"github.com/globalsign/mgo"
)

var db *mgo.Database

// InitializeDB initates database connection and connects to specified database
func InitializeDB(dbName string) error {
	var dbSession, err = mgo.Dial("localhost")
	if err != nil {
		return err
	}
	db = dbSession.DB(dbName)

	db.C(httpsessionCollection).EnsureIndex(mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: 15 * time.Minute})

	db.C(httpsessionCollection).EnsureIndex(mgo.Index{
		Key: []string{"key"}})

	db.C(autoMessageCollectionName).EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	db.C(channelBansCollectionName).EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	db.C(channelInfoCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	db.C(gamesCollection).EnsureIndex(mgo.Index{
		Key: []string{"gameid"}})

	return nil
}
