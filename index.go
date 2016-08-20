package main

import (
	"log"
	"time"

	"github.com/khades/servbot/bot"
)

func main() {
	log.Println("Starting...")
	ticker := time.NewTicker(time.Second * 15)
	go func() {
		for {
			tick := <-ticker.C
			log.Print(tick)
			bot.IrcClientInstance.SendModsCommand()
		}
	}()
	bot.Start()
}
