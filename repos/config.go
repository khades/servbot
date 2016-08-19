package repos

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/khades/servbot/models"
)

func readConfig() models.ConfigModel {
	var configfile string
	flag.StringVar(&configfile, "config", "config.json", "defines configuration file for application")

	file, err := ioutil.ReadFile(configfile)

	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}
	var config models.ConfigModel
	error := json.Unmarshal(file, &config)
	if error != nil {
		log.Fatal("json read error", error)
	}
	if config.OauthKey == "" || config.BotUserName == "" || config.DbName == "" {
		log.Fatal("config.json has invalid format")
	}
	return config
}

// Config is object about config
var Config = readConfig()
