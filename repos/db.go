package repos

import "gopkg.in/mgo.v2"

var dbSession, err = mgo.Dial("localhost")

// Db is database connection object
var Db = dbSession.DB(Config.DbName)
