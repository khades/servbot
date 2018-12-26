package twitchIRCCTasks

import (
	"github.com/khades/servbot/twitchIRC"
	"time"
)

func Run(twitchIRCClient *twitchIRC.Client) {
	ircClientTicker := time.NewTicker(time.Second * 3)

	go func() {
		for {

			<-ircClientTicker.C
			twitchIRCClient.SendMessages(3)
		}
	}()

	modTicker := time.NewTicker(time.Second * 10)

	go func() {
		for {
			<-modTicker.C
			twitchIRCClient.SendModsCommand()
		}
	}()

}
