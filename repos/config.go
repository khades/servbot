package repos

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/khades/servbot/models"
	"gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

var configCollection = "config"

// ReadConfigFromFile returns config object parsed from config.json
func ReadConfigFromFile() models.Config {
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

// ReadConfigFromDatabase returns config object parsed from mongodb
func ReadConfigFromDatabase() (models.Config, error) {
	var result models.Config
	error := db.C(configCollection).Find(bson.M{"entity": "config"}).One(&result)
	return result, error
}

// SaveConfigToDatabase saves current config object to database
func SaveConfigToDatabase() {
	db.C(configCollection).Upsert(bson.M{"entity": "config"}, bson.M{"$set": Config})

}

// Config represents config file as object, it is populated from config.json at bot initialization phase
var Config models.Config
