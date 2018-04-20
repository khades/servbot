package main

import (
	"flag"
	"log"

	"github.com/khades/servbot/pubsub"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	dbName := flag.String("db", "servbot", "mongo database name")
	// Initializing database
	dbErr := repos.InitializeDB(*dbName)
	if dbErr != nil {
		log.Fatalf("Database Conenction Error: " + dbErr.Error())
	}
	localConfig, configError := repos.ReadConfigFromDatabase()

	if configError != nil {
		log.Fatalf("Reading config from database failed: %s", configError)
	}

	repos.Config = localConfig
	pubsub.RunPubSub()
}
