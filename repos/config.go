package repos

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/asaskevich/govalidator.v4"
	"github.com/khades/servbot/models"
)

func readConfig() models.Config {
	var configfile string
	flag.StringVar(&configfile, "config", "config.json", "defines configuration file for application")

	file, err := ioutil.ReadFile(configfile)

	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}
	var config models.Config
	error := json.Unmarshal(file, &config)
	if error != nil {
		log.Fatal("json read error", error)
	}
	result, err := govalidator.ValidateStruct(config)
	if result == false || err != nil {
		log.Fatal("config.json has invalid format: ", err)

	}

	return config
}

// Config represents config file as object
var Config = readConfig()
