package main

import (
	"encoding/gob"
	"log"
	"sync"
	"time"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/services"
)

func main() {
	var wg sync.WaitGroup

	gob.Register(&models.HTTPSession{})
	log.Println("Starting...")
	ticker := time.NewTicker(time.Second * 15)

	go func(wg *sync.WaitGroup) {
		for {
			<-ticker.C
			wg.Add(1)
			bot.IrcClientInstance.SendModsCommand()
			wg.Done()
		}
	}(&wg)

	minuteTicker := time.NewTicker(time.Minute)
	services.CheckStreamStatus()

	go func(wg *sync.WaitGroup) {
		for {
			<-minuteTicker.C
			wg.Add(1)
			services.CheckStreamStatus()
			wg.Done()
		}
	}(&wg)

	// go func(wg *sync.WaitGroup) {
	// 	//defer wg.Done()
	// 	wg.Add(1)
	// 	httpbackend.Start()

	// }(&wg)

	// go func(wg *sync.WaitGroup) {
	// 	wg.Add(1)
	// 	//bot.Start()
	// 	wg.Done()
	// }(&wg)
	bot.Start()

	log.Println("Quitting...")
	wg.Wait()

}
