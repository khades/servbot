package repos

import "gopkg.in/mgo.v2"

var db *mgo.Database 

// InitializeDB initates database connection and connects to specified database
func InitializeDB(dbName string) error {
	var dbSession, err = mgo.Dial("localhost")
	if err != nil {
		return err
	}
	db = dbSession.DB(dbName)
	return nil
}