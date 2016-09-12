package bot

import "github.com/khades/servbot/ircClient"

// IrcClientInstance is concrete irc client we will work with
var IrcClientInstance = ircClient.IrcClient{Ready: false}
