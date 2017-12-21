package main

import (
	"encoding/gob"
	"log"
	"sync"
	"time"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/httpbackend"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/khades/servbot/services"
)

func main() {
	var wg sync.WaitGroup
	result, resultError := repos.GetUsersID(&repos.Config.Channels)
	if resultError != nil {
		log.Printf("INITIALISATION ERROR: ", resultError)
		return
	}
	for _, value := range *result {
		repos.PushCommandsForChannel(&value)
	}

	repos.CreateChannels()
	wg.Add(1)
	go func() {
		services.CheckTwitchDJTrack()
		services.CheckStreamStatus()
		// 	services.CheckDubTrack()
	}()

	gob.Register(&models.HTTPSession{})
	log.Println("Starting...")

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
	// go func(wg *sync.WaitGroup) {
	// 	for {
	// 		<-thirtyTicker.C
	// 		wg.Add(1)
	// 		services.CheckDubTrack()
	// 		wg.Done()
	// 	}
	// }(&wg)
	twentyticker := time.NewTicker(time.Second * 20)

	go func() {
		for {
			<-twentyticker.C
			eventbus.EventBus.Trigger("ping")
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
			services.CheckStreamStatus()
			wg.Done()
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		httpbackend.Start()
		wg.Done()
	}(&wg)
	followerTicker := time.NewTicker(time.Second * 30)

	go func(wg *sync.WaitGroup) {
		for {
			<-followerTicker.C
			wg.Add(1)
			services.CheckChannelsFollowers()
			wg.Done()
		}
	}(&wg)
	go func(wg *sync.WaitGroup) {
		bot.Start()
		wg.Done()
	}(&wg)

	wg.Wait()
	log.Println("Quitting...")

}
