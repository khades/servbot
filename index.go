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
	repos.GetUsersID(&repos.Config.Channels)
	repos.CreateChannels()
	wg.Add(1)
	go func() {
		services.CheckTwitchDJTrack()
		services.CheckStreamStatus()
		// 	services.CheckDubTrack()
	}()

	gob.Register(&models.HTTPSession{})
	log.Println("Starting...")
	ticker := time.NewTicker(time.Second * 15)

	go func(wg *sync.WaitGroup) {
		for {
			<-ticker.C
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

	go func(wg *sync.WaitGroup) {
		bot.Start()
		wg.Done()
	}(&wg)

	wg.Wait()
	log.Println("Quitting...")

}
