package main

import (
	"log"
	"time"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/services"
)

func main() {
	log.Println("Starting...")
	services.CheckStreamStatus()

	ticker := time.NewTicker(time.Second * 15)
	go func() {
		for {
			tick := <-ticker.C
			log.Print(tick)
			bot.IrcClientInstance.SendModsCommand()
		}
	}()

	minuteTicker := time.NewTicker(time.Minute)
	go func() {
		for {
			tick := <-minuteTicker.C
			log.Print(tick)
			services.CheckStreamStatus()
		}
	}()
	bot.Start()
}
