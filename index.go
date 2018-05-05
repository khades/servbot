package main

import (
	"encoding/gob"
	"flag"
	"sync"
	"time"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/pubsub"
	"github.com/sirupsen/logrus"

	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/httpbackend"
	"github.com/khades/servbot/repos"
	"github.com/khades/servbot/services"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	logger := logrus.WithFields(logrus.Fields{"package": "main"})
	logger.Info("Starting")
	convertConfig := flag.Bool("convertconfig", false, "forces importing config file to database")
	dbName := flag.String("db", "servbot", "mongo database name")
	flag.Parse()
	logger.Infof("Database name: %s", *dbName)
	// Initializing database
	dbErr := repos.InitializeDB(*dbName)
	if dbErr != nil {
		logger.Fatal("Database Conenction Error: " + dbErr.Error())
	}

	if *convertConfig == true {
		logrus.SetLevel(logrus.DebugLevel)

		logger.Info("Running configuration importer.")
		repos.Config = repos.ReadConfigFromFile()
		users, usersError := repos.GetUsersID(repos.Config.Channels)
		if usersError != nil {
			logger.Fatalf("User conversion error: %s", usersError.Error())
		}
		channelIDs := []string{}
		for _, value := range *users {
			channelIDs = append(channelIDs, value)
			repos.EnableChannel(&value)
			lang := "ru"
			repos.SetChannelLang(&value, &lang)
		}
		repos.SaveConfigToDatabase()

		logger.Info("Configuration import successed.")

		return
	}

	// Database initialisation and preprocessing

	// Reading config from database
	localConfig, configError := repos.ReadConfigFromDatabase()

	if configError != nil {
		logger.Fatalf("Reading config from database failed: %s", configError)
	}

	repos.Config = localConfig
	if repos.Config.Debug == true {
		logrus.SetLevel(logrus.DebugLevel)
	}

	repos.PreprocessChannels()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		services.CheckTwitchDJTrack()
		services.CheckStreamStatuses()
		// 	services.CheckDubTrack()
	}()

	gob.Register(&models.HTTPSession{})
	logger.Info("Starting...")
	//services.CheckChannelsFollowers()
	ircClientTicker := time.NewTicker(time.Second * 3)

	go func(wg *sync.WaitGroup) {
		for {
			wg.Add(1)
			<-ircClientTicker.C
			bot.IrcClientInstance.SendMessages(3)
			wg.Done()

		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		for {
			wg.Add(1)
			pubsub.RunPubSub()
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
	modTicker := time.NewTicker(time.Second * 30)

	go func(wg *sync.WaitGroup) {
		for {
			<-modTicker.C
			wg.Add(1)
			bot.IrcClientInstance.SendModsCommand()
			services.SendAutoMessages()
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

	// webhookTimer := time.NewTicker(time.Second * 300)
	// repos.CheckAndSubscribeToWebhooks(time.Second * 300)
	//repos.SubChannelToFollowerHooks("40635840")
	// go func() {
	// 	for {
	// 		<-webhookTimer.C
	// 		repos.CheckAndSubscribeToWebhooks(time.Second * 300)
	// 	}
	// }()
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
	// followerTicker := time.NewTicker(time.Second * 30)

	// go func(wg *sync.WaitGroup) {
	// 	for {
	// 		<-followerTicker.C
	// 		wg.Add(1)
	// 		services.CheckChannelsFollowers()
	// 		wg.Done()
	// 	}
	// }(&wg)
	go func(wg *sync.WaitGroup) {
		bot.Start()
		wg.Done()
	}(&wg)

	wg.Wait()
	logger.Info("Quitting...")
	// Kseyko = PIDR
}
