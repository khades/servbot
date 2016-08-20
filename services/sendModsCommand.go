package services

// SendModsCommand
func SendModsCommand() {
        for _, value := range repos.Config.Channels {
	    	bot.IrcClient.SendRaw("#PRIVMSG "+value+" .mods")
		}
}