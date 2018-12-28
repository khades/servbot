package twitchIRCCTasks

import (
	"github.com/khades/servbot/twitchIRC"
	"time"
)

func Run(twitchIRCClient *twitchIRC.Client) {
	ircClientTicker := time.NewTicker(time.Second * 3)

	go func() {
		for range ircClientTicker.C{
			twitchIRCClient.SendMessages(3)
		}
	}()

	modTicker := time.NewTicker(time.Second * 10)

	go func() {
		for range modTicker.C {
			twitchIRCClient.SendModsCommand()
		}
	}()

}
