package main

import (
	"encoding/gob"
	"log"
	"time"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/services"
)

func main() {
	gob.Register(&models.HTTPSession{})
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
