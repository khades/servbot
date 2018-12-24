package main

import (
	"encoding/gob"
	"flag"
	"sync"
	"time"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/httpbackend"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/pubsub"
	"github.com/khades/servbot/repos"
	"github.com/khades/servbot/services"
	"github.com/sirupsen/logrus"
)

func main2() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	logger := logrus.WithFields(logrus.Fields{"package": "main"})
	logger.Info("Starting")

	dbName := flag.String("db", "servbot", "mongo database name")
	flag.Parse()
	logger.Infof("Database name: %s", *dbName)
	// Initializing database
	var dbSession, err = mgo.Dial("localhost")
	if err != nil {
		logger.Fatal("Database Conenction Error: " + err.Error())
	}
	db := dbSession.DB(dbName)

	// Reading config from database
	config, configError := config.Init(db)

	if configError != nil {
		logger.Fatalf("Reading config from database failed: %s", configError)
	}

	if config.Debug == true {
		logrus.SetLevel(logrus.DebugLevel)
	}

	twitchAPIService := twitchAPI.Init(config)
	userResolveService = userResolve.Init(db, twitchAPIService)
	repos.PreprocessChannels()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		services.CheckTwitchDJTrack()
		services.CheckStreamStatuses()
	}()

	gob.Register(&models.HTTPSession{})
	logger.Info("Starting...")

	ircClientTicker := time.NewTicker(time.Second * 3)

	go func(wg *sync.WaitGroup) {
		for {
			wg.Add(1)
			<-ircClientTicker.C
			bot.IrcClientInstance.SendMessages(3)
			wg.Done()

		}
	}(&wg)
	modTicker := time.NewTicker(time.Second * 10)

	go func(wg *sync.WaitGroup) {
		for {
			<-modTicker.C
			wg.Add(1)
			bot.IrcClientInstance.SendModsCommand()
			services.SendAutoMessages()
			wg.Done()
		}
	}(&wg)
	go func(wg *sync.WaitGroup) {
		for {
			wg.Add(1)
			pubsub.Run()
			wg.Done()

		}
	}(&wg)
	gamesCheckerTicker := time.NewTicker(time.Second * 30)
	go func(wg *sync.WaitGroup) {
		for {
			<-gamesCheckerTicker.C
			wg.Add(1)
			services.GetTwitchGames()
			wg.Done()
		}
	}(&wg)


	thirtyTicker := time.NewTicker(time.Second * 30)
	go func(wg *sync.WaitGroup) {
		for {
			<-thirtyTicker.C
			wg.Add(1)
			services.CheckTwitchDJTrack()
			wg.Done()
		}
	}(&wg)

	subTrainNotificationTicker := time.NewTicker(time.Second * 5)
	go func(wg *sync.WaitGroup) {
		for {
			<-subTrainNotificationTicker.C
			wg.Add(1)
			services.SendSubTrainNotification()
			wg.Done()
		}
	}(&wg)

	subTrainTimeoutTicker := time.NewTicker(time.Second * 5)
	go func(wg *sync.WaitGroup) {
		for {
			<-subTrainTimeoutTicker.C
			wg.Add(1)
			services.SendSubTrainTimeoutMessage()
			wg.Done()
		}
	}(&wg)

	pingticker := time.NewTicker(time.Second * 30)

	go func() {
		for {
			<-pingticker.C
			eventbus.EventBus.Publish("ping", "ping")
		}
	}()

	webhookTimer := time.NewTicker(time.Minute * 15)
	repos.CheckAndSubscribeToWebhooks(time.Minute * 15)
	go func() {
		for {
			<-webhookTimer.C
			repos.CheckAndSubscribeToWebhooks(time.Minute * 15)
		}
	}()

	vkTimer := time.NewTicker(time.Second * 60)

	go func() {
		for {
			<-vkTimer.C
			services.CheckVK()
		}
	}()
	minuteTicker := time.NewTicker(time.Minute)

	go func(wg *sync.WaitGroup) {
		for {
			<-minuteTicker.C
			wg.Add(1)
			services.CheckStreamStatuses()
			wg.Done()
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		httpbackend.Start()
		wg.Done()
	}(&wg)

	followerTicker := time.NewTicker(time.Second * 20)

	go func(wg *sync.WaitGroup) {
		for {
			<-followerTicker.C
			wg.Add(1)
			services.AnnounceFollowers()
			wg.Done()
		}
	}(&wg)
	go func(wg *sync.WaitGroup) {
		bot.Start()
		wg.Done()
	}(&wg)

	wg.Wait()
	logger.Info("Quitting...")
	// Kseyko = PIDR
}
