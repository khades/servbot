package main

import (
	"flag"
	"sync"
	"time"

	"github.com/khades/servbot/donationSourceTasks"
	"github.com/khades/servbot/yandexOAuth"

	"github.com/khades/servbot/donationAPI"
	"github.com/khades/servbot/donationSource"
	"github.com/khades/servbot/donationSourceAPI"
	"github.com/khades/servbot/event"

	"github.com/khades/servbot/balance"
	"github.com/khades/servbot/currencyConverter"
	"github.com/khades/servbot/donation"
	"github.com/khades/servbot/metrics"

	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/autoMessageTasks"
	"github.com/khades/servbot/configTasks"
	"github.com/khades/servbot/followersToGreetTasks"
	"github.com/khades/servbot/pubsub"
	"github.com/khades/servbot/streamStatusTasks"
	"github.com/khades/servbot/subtrainTasks"
	"github.com/khades/servbot/twitchIRC"
	"github.com/khades/servbot/twitchIRCCTasks"
	"github.com/khades/servbot/twitchIRCHandler"
	"github.com/khades/servbot/videoLibraryAPI"
	"github.com/khades/servbot/vkGroupTasks"
	"github.com/khades/servbot/webhookTasks"

	"github.com/khades/servbot/autoMessageAPI"
	"github.com/khades/servbot/channelBansAPI"
	"github.com/khades/servbot/channelLogsAPI"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/songRequestAPI"
	"github.com/khades/servbot/subAlertAPI"
	"github.com/khades/servbot/subdayAPI"
	"github.com/khades/servbot/subscriptionInfoAPI"
	"github.com/khades/servbot/subtrainAPI"
	"github.com/khades/servbot/templateAPI"
	"github.com/khades/servbot/vkGroupAPI"
	"github.com/khades/servbot/webhookAPI"

	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/subscriptionInfo"
	"github.com/khades/servbot/webhook"

	"github.com/khades/servbot/followers"

	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/songRequest"

	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/subday"
	"github.com/khades/servbot/videoLibrary"
	"github.com/khades/servbot/youtubeAPI"

	"github.com/khades/servbot/channelBans"
	"github.com/khades/servbot/channelLogs"

	"github.com/khades/servbot/gameResolve"
	"github.com/khades/servbot/streamStatus"
	"github.com/khades/servbot/template"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchAPI"
	"github.com/khades/servbot/userResolve"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

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
	db := dbSession.DB(*dbName)

	// Reading config from database
	config, configError := config.Init(db)

	if config.Debug == true {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if configError != nil {
		logger.Fatalf("Reading config from database failed: %s", configError)
	}
	// Creating waitgroup for timed services
	var wg sync.WaitGroup

	var metrics = metrics.Init()

	var eventBus = EventBus.New()

	// Creating ticker for websocket ping events
	pingticker := time.NewTicker(time.Second * 30)

	go func() {
		for {
			<-pingticker.C
			eventBus.Publish("ping", "ping")
		}
	}()

	currencyConverterService := currencyConverter.Init()

	balanceService := balance.Init(db)

	eventService := event.Init(db)

	donationService := donation.Init(db, currencyConverterService, balanceService, eventService)

	// Creating twitchAPI service
	twitchAPIClient := twitchAPI.Init(
		config,
		metrics)
	configTasks.Run(twitchAPIClient, config)

	// Creating youtubeAPI service
	youtubeAPIClient := youtubeAPI.Init(
		config)

	// Service
	httpSessionService := httpSession.Init(
		db,
		twitchAPIClient)

	// Creating user to userID resolution service
	userResolveService := userResolve.Init(
		db,
		twitchAPIClient)

	// Creating channelInfo super-object service
	channelInfoService := channelInfo.Init(
		db,
		config,
		userResolveService,
		metrics)

	// HttpAPI
	httpAPIService := httpAPI.Init(
		config,
		httpSessionService,
		channelInfoService,
		eventBus,
		metrics)

	subtrainAPI.Init(
		httpAPIService,
		channelInfoService)

	vkGroupAPI.Init(
		httpAPIService,
		channelInfoService)

	donationAPI.Init(
		httpAPIService,
		donationService)

	donationSourceService := donationSource.Init(db)
	donationSourceAPI.Init(httpAPIService, donationSourceService)
	yandexOAuth.Init(httpAPIService, config, donationSourceService)
	donationSourceTasks.Run(donationService, donationSourceService)
	// Creating gameID to game resolution service
	gameResolveService, _ := gameResolve.Init(
		db,
		twitchAPIClient,
		channelInfoService)

	// Creating twitchAPI stream status checker service
	streamStatusService := streamStatus.Init(
		config,
		channelInfoService,
		gameResolveService,
		twitchAPIClient)

	streamStatusTasks.Run(
		streamStatusService)

	// Creating Template service
	templateService := template.Init(
		db,
		channelInfoService)

	// Initialising TemplateAPI
	templateAPI.Init(
		httpAPIService,
		templateService)

	// Creating ChannelBans service
	channelBansService := channelBans.Init(
		db)

	// Initialising ChannelBansAPI
	channelBansAPI.Init(
		httpAPIService,
		channelBansService)

	// Creating Service
	channelLogsService := channelLogs.Init(
		db,
		channelBansService,
		userResolveService)

	// Attaching api methods
	channelLogsAPI.Init(
		httpAPIService,
		channelLogsService)

	// Creating Service
	subdayService := subday.Init(
		db,
		channelInfoService)

	// Registering subday to httpAPI
	subdayAPI.Init(
		httpAPIService,
		subdayService)

	// Creating Service
	subAlertService := subAlert.Init(
		db)

	// API
	subAlertAPI.Init(
		httpAPIService,
		subAlertService)

	// Creating videoLibraryService
	videoLibraryService := videoLibrary.Init(
		db)

	// Creating songRequestService
	songRequestService := songRequest.Init(
		db,
		youtubeAPIClient,
		channelInfoService,
		videoLibraryService,
		eventBus)

	songRequestAPI.Init(
		httpAPIService,
		songRequestService)

	// Service
	followersToGreetService := followersToGreet.Init(
		db)

	// Service
	followersService := followers.Init(
		db,
		twitchAPIClient,
		followersToGreetService)

	// Running PubSub
	pubsub := pubsub.Init(
		channelInfoService,
		config,
		channelLogsService,
		&wg)

	// Running Webhooks
	webHookService := webhook.Init(
		db,
		channelInfoService,
		twitchAPIClient)

	webhookTasks.Run(
		webHookService)

	// Constructing webhooksAPI
	webhookAPI.Init(
		httpAPIService,
		webHookService,
		streamStatusService,
		followersService)

	// Automessage service
	autoMessageService := autoMessage.Init(db)

	// Its API
	autoMessageAPI.Init(httpAPIService, autoMessageService)

	// SubscriptionInfo service
	subscriptionInfoService := subscriptionInfo.Init(db)

	subscriptionInfoAPI.Init(
		httpAPIService,
		subscriptionInfoService,
		eventBus,
	)

	// TwitchIRCHandler
	twitchIRCHandler := twitchIRCHandler.Init(
		subdayService,
		channelInfoService,
		subAlertService,
		channelLogsService,
		autoMessageService,
		userResolveService,
		subscriptionInfoService,
		templateService,
		followersService,
		songRequestService,
		eventBus,
		pubsub,
		eventService,
		balanceService,
	)

	// TwitchBot
	twitchIRCClient := twitchIRC.Init(
		config,
		channelInfoService,
		twitchIRCHandler.Handle,
		metrics,
		&wg,
	)
	twitchIRCCTasks.Run(
		twitchIRCClient)

	// Initialising videoLibraryAPI
	videoLibraryAPI.Init(
		httpAPIService,
		videoLibraryService,
		songRequestService,
		twitchIRCClient,
		eventBus)
	// FollowerAnnouncer
	followersToGreetTasks.Run(
		channelInfoService,
		followersToGreetService,
		subAlertService,
		userResolveService,
		twitchIRCClient)

	// AutomessageAnnouncer
	autoMessageTasks.Run(
		channelInfoService,
		autoMessageService,
		twitchIRCClient)

	// SubtrainAnnouncer
	subtrainTasks.Run(
		channelInfoService,
		twitchIRCClient,
		eventBus)

	vkGroupTasks.Run(
		config,
		channelInfoService,
		twitchIRCClient)

	httpAPIService.Serve(&wg)
	wg.Wait()

}
