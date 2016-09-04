package bot

import "github.com/khades/servbot/ircClient"

// IrcClientInstance is concrete irc client we work with
var IrcClientInstance = ircClient.IrcClient{Ready: false}
