package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type kv struct {
	Key   int64
	Value int64
}

func main() {

	var db *mgo.Database
	testCollection := "test"
	// InitializeDB initates database connection and connects to specified database

	var dbSession, err = mgo.Dial("localhost")
	if err != nil {
		return
	}

	//db.C(testCollection).DropCollection()
	db = dbSession.DB("test")
	for i := 1; i <= 10000; i++ {
		bulk := db.C(testCollection).Bulk() // Getting error. (writeConcern: { getLastError: 1 })
		bulk.Unordered()
		for s := 1; s <= 1000; s++ {
			bulk.Insert(
				bson.M{"value": i*1000 + s, "createdat": time.Now()})
		}
		bulk.Run()

	}
	db.C(testCollection).EnsureIndex(mgo.Index{
		Key: []string{"key"}})

	db.C(testCollection).EnsureIndex(mgo.Index{
		Key: []string{"value"}})

	db.C(testCollection).EnsureIndex(mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: 20 * time.Minute})
	values := []int64{}
	for s := 1; s <= 50; s++ {
		for i := 1; i <= 50; i++ {
			values = append(values, rand.Int63n(10000000))
		}
		result := []kv{}
		db.C(testCollection).Find(bson.M{"key": bson.M{"$in": values}}).All(&result)
		log.Println(result)
	}

}
