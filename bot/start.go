package bot

import (
	"log"
	"net"

	"github.com/belak/irc"
	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/repos"
)

type chatClient struct {
	Client *irc.Client
	Ready  bool
}

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
	chatClient := irc.NewClient(conn, config)
	log.Println("Bot is starting...")

	clientError := chatClient.Run()
	log.Fatalln(clientError)
	IrcClientInstance = ircClient.IrcClient{Ready: false}
	conn.Close()
}
