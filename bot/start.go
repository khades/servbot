package bot

import (
	"net"

	"github.com/khades/servbot/ircClient"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v2"
)

type chatClient struct {
	Client *irc.Client
	Ready  bool
}

// Start function dials up connection for chathandler
func Start() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "bot",
		"feature": "bot",
		"action":  "Start"})
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		logger.Fatalln(err)
	}

	config := irc.ClientConfig{
		Nick:    repos.Config.BotUserName,
		Pass:    repos.Config.OauthKey,
		User:    repos.Config.BotUserName,
		Name:    repos.Config.BotUserName,
		Handler: chatHandler}
	chatClient := irc.NewClient(conn, config)
	logger.Info("Bot is starting")

	clientError := chatClient.Run()
	logger.Info(clientError)
	logger.Fatal("Bot died")
	IrcClientInstance = &ircClient.IrcClient{Ready: false, MessageQueue: []string{}}
	conn.Close()
}
