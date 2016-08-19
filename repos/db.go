package repos

import "gopkg.in/mgo.v2"

var session, err = mgo.Dial("localhost")

// Db is database object
// TODO name the db in config
var Db = session.DB(Config.DbName)
