package videoLibrary

import "github.com/globalsign/mgo"

const videolibraryCollection = "videolibrary"

// Init initalises mongo collection, creates indexes and returns service
func Init(db *mgo.Database) *Service {
	collection := db.C(videolibraryCollection)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"videoid"}})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"-_id"}})

	collection.EnsureIndex(mgo.Index{
		Key: []string{"tags.tag"}})

	return &Service{
		collection:            collection,
		count:                 -1,
		bannedCountPerChannel: make(map[string]int),
	}
}
