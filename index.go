package main

import (
	"log"
	"time"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/services"
)

func main() {
	log.Println("Starting...")
	ticker := time.NewTicker(time.Second * 15)
	go func() {
		for {
			<-ticker.C
			bot.IrcClientInstance.SendModsCommand()
		}
	}()

	minuteTicker := time.NewTicker(time.Minute)
	services.CheckStreamStatus()

	go func() {
		for {
			<-minuteTicker.C
			services.CheckStreamStatus()
		}
	}()
	// go httpbackend.Start()
	bot.Start()
}
