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

	db.C(followerCursorsCollectionName).EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	db.C(subAlertCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	db.C(songRequestCollectionName).EnsureIndex(mgo.Index{
		Key: []string{"channelid"}})

	db.C(videolibraryCollection).EnsureIndex(mgo.Index{
		Key: []string{"videoid"}})

	db.C(videolibraryCollection).EnsureIndex(mgo.Index{
		Key: []string{"-_id"}})

	db.C(videolibraryCollection).EnsureIndex(mgo.Index{
		Key: []string{"tags.tag"}})

	db.C(gamesCollection).EnsureIndex(mgo.Index{
		Key: []string{"gameid"}})

	db.C(usernameCacheCollection).EnsureIndex(mgo.Index{
		Key: []string{"id"}})

	db.C(usernameCacheCollection).EnsureIndex(mgo.Index{
		Key: []string{"userid"}})

	db.C(subscriptionInfoCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})

	db.C(templateCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "commandname"}})

	db.C(followersListСollectionName).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})

	db.C("webhooklibrary").EnsureIndex(mgo.Index{
		Key: []string{"channelid", "topic"}})

	db.C(messageLogsCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "userid"}})
	db.C(messageLogsCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "-lastupdate"}})
	db.C(messageLogsCollection).EnsureIndex(mgo.Index{
		Key: []string{"channelid", "$text:knownnicknames"}})

	db.C(usernameCacheCollection).EnsureIndex(mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: 60 * 12 * time.Minute})

	return nil
}
