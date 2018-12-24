package main

import (
	"flag"
	evbus "github.com/asaskevich/EventBus"
	"github.com/khades/servbot/autoMessageAPI"
	"github.com/khades/servbot/autoMessageAnnounce"
	"github.com/khades/servbot/channelBansAPI"
	"github.com/khades/servbot/channelLogsAPI"
	"github.com/khades/servbot/followersAnnounce"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/songRequestAPI"
	"github.com/khades/servbot/subAlertAPI"
	"github.com/khades/servbot/subdayAPI"
	"github.com/khades/servbot/subscriptionInfoAPI"
	"github.com/khades/servbot/subtrainAPI"
	"github.com/khades/servbot/subtrainAnnounce"
	"github.com/khades/servbot/templateAPI"
	"github.com/khades/servbot/videoLibraryAPI"
	"github.com/khades/servbot/vkGroupAPI"
	"github.com/khades/servbot/vkGroupAnnounce"
	"github.com/khades/servbot/webhookAPI"
	"sync"
	"time"

	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/subscriptionInfo"
	"github.com/khades/servbot/twitchIRCClient"
	"github.com/khades/servbot/twitchIRCHandler"
	"github.com/khades/servbot/webhook"

	"github.com/khades/servbot/followers"

	"github.com/khades/servbot/followersToGreet"
	"github.com/khades/servbot/songRequest"

	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/subday"
	"github.com/khades/servbot/videoLibrary"
	"github.com/khades/servbot/youtubeAPIClient"

	"github.com/khades/servbot/channelBans"
	"github.com/khades/servbot/channelLogs"

	"github.com/khades/servbot/gameResolve"
	"github.com/khades/servbot/streamStatus"
	"github.com/khades/servbot/template"

	"github.com/globalsign/mgo"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchAPIClient"
	"github.com/khades/servbot/userResolve"
	"github.com/sirupsen/logrus"
)

func main() {
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
	db := dbSession.DB(*dbName)

	// Reading config from database
	config, configError := config.Init(db)

	if configError != nil {
		logger.Fatalf("Reading config from database failed: %s", configError)
	}
	// Creating waitgroup for timed services
	var wg sync.WaitGroup

	var eventBus = evbus.New()

	pingticker := time.NewTicker(time.Second * 30)

	go func() {
		for {
			<-pingticker.C
			eventBus.Publish("ping", "ping")
		}
	}()

	tickers := []*time.Ticker{}

	//TODO Decoumple timers from services
	if config.Debug == true {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// Creating twitchAPIClient service
	twitchAPIClient := twitchAPIClient.Init(
		config,
	)

	// Creating youtubeAPIClient service
	youtubeAPIClient := youtubeAPIClient.Init(
		config,
	)

	// Service
	httpSessionService := httpSession.Init(
		db,
		twitchAPIClient,
	)

	// Creating user to userID resolution service
	userResolveService := userResolve.Init(
		db,
		twitchAPIClient,
	)

	// Creating channelInfo super-object service
	channelInfoService := channelInfo.Init(
		db,
		config,
		userResolveService,
	)

	// HttpAPI
	httpAPIService := httpAPI.Init(
		config,
		httpSessionService,
		channelInfoService,
		eventBus,
	)

	subtrainAPI.Init(httpAPIService, channelInfoService)
	vkGroupAPI.Init(httpAPIService, channelInfoService)

	// Creating gameID to game resolution service
	gameResolveService, _ := gameResolve.Init(
		db,
		twitchAPIClient,
		channelInfoService,
		&wg,
	)

	// Creating twitchAPIClient stream status checker service
	streamStatusService, streamStatusTicker := streamStatus.Init(
		config,
		channelInfoService,
		gameResolveService,
		twitchAPIClient,
		&wg,
	)
	tickers = append(tickers, streamStatusTicker)

	// Creating Template service
	templateService :=template.Init(
		db,
		channelInfoService,
	)

	// Initialising TemplateAPI
	templateAPI.Init(httpAPIService, templateService)

	// Creating ChannelBans service
	channelBansService := channelBans.Init(
		db,
	)

	// Initialising ChannelBansAPI
	channelBansAPI.Init(httpAPIService, channelBansService)

	// Creating Service
	channelLogsService := channelLogs.Init(
		db,
		channelBansService,
		userResolveService,
	)

	// Attaching api methods
	channelLogsAPI.Init(httpAPIService, channelLogsService)

	// Creating Service
	subdayService := subday.Init(
		db,
		channelInfoService,
	)

	// Registering subday to httpAPI
	subdayAPI.Init(httpAPIService,subdayService)

	// Creating Service
	subAlertService := subAlert.Init(
		db,
	)
	// API
	subAlertAPI.Init(httpAPIService, subAlertService)

	// Creating videoLibraryService
	videoLibraryService :=
		videoLibrary.Init(
			db,
		)

	// Creating songRequestService
	songRequestService := songRequest.Init(
		db,
		youtubeAPIClient,
		channelInfoService,
		videoLibraryService,
		eventBus,
	)
	songRequestAPI.Init(httpAPIService, songRequestService)

	// Service
	followersToGreetService := followersToGreet.Init(
		db,
	)

	// Service
	followersService :=	followers.Init(
		db,
		twitchAPIClient,
		followersToGreetService,
	)

	// Running PubSub
	// TODO: it runs once!
	// pubsub.RunPubSub(channelInfoService, config, channelLogsService)

	// Running Webhooks
	webHookService, webHookTicker := webhook.Init(
		db,
		channelInfoService,
		twitchAPIClient,
		&wg,
	)
	tickers = append(tickers, webHookTicker)

	// Constructing webhooksAPI
	webhookAPI.Init(httpAPIService, webHookService, streamStatusService, followersService )

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
	)

	// TwitchBot
	twitchIRCClient, twitchIRCTicker, twitchIRCModTicker := twitchIRCClient.Init(
		config,
		channelInfoService,
		twitchIRCHandler.Handle,
		&wg,
	)
	tickers = append(tickers, twitchIRCTicker)
	tickers = append(tickers, twitchIRCModTicker)

	// Initialising videoLibraryAPI
	videoLibraryAPI.Init(
		httpAPIService,
		videoLibraryService,
		songRequestService,
		twitchIRCClient,
		eventBus,
	)
	// FollowerAnnouncer
	followersTicker := followersAnnounce.Init(
		channelInfoService,
		followersToGreetService,
		subAlertService,
		userResolveService,
		twitchIRCClient,
		&wg,
	)
	tickers = append(tickers, followersTicker)

	// AutomessageAnnouncer
	automessageAnnounceTicker := autoMessageAnnounce.Init(
		channelInfoService,
		autoMessageService,
		twitchIRCClient,
		&wg,
	)
	tickers = append(tickers, automessageAnnounceTicker)

	// SubtrainAnnouncer
	subtrainAnnounceTicker := subtrainAnnounce.Init(
		channelInfoService,
		twitchIRCClient,
		eventBus,
		&wg,
	)
	tickers = append(tickers, subtrainAnnounceTicker)

	vkAnnounceTicker := vkGroupAnnounce.Init(
		config,
		channelInfoService,
		twitchIRCClient,
		&wg)
	tickers = append(tickers, vkAnnounceTicker)

	httpAPIService.Serve()
	wg.Wait()

}
