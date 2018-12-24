package config

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v4"
)

var configCollection = "config"

// Init returns config object parsed from mongodb
func Init(db *mgo.Database) (*Config, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "config",
		"feature": "init",
		"action":  "init"})
	var config Config
	config.collection = db.C(configCollection)

	error := config.collection.Find(bson.M{"entity": "config"}).One(&config)

	if error != nil {
		return nil, error
	}

	validated, validationError := govalidator.ValidateStruct(config)

	if validated == true {
		logger.Infof("Config parsing error: %+v", validationError)
		return nil, error
	}

	return &config, error
}
