package bot

import (
	"log"
	"net"

	"github.com/belak/irc"
	"github.com/khades/servbot/repos"
)

// Start starts the bot
func Start() {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Fatalln(err)
	}

	config := irc.ClientConfig{
		Nick:    repos.Config.BotUserName,
		Pass:    repos.Config.OauthKey,
		User:    repos.Config.BotUserName,
		Name:    repos.Config.BotUserName,
		Handler: chatHandler}

	clientError := irc.NewClient(conn, config).Run()
	log.Fatalln(clientError)
	conn.Close()
}
