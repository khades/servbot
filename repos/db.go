package repos

import "gopkg.in/mgo.v2"

var session, err = mgo.Dial("localhost")

// Db is database connection object
var Db = session.DB(Config.DbName)
