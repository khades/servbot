package config

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v4"
)

const collectionName = "config"

// Init returns config object parsed from mongodb
func Init(db *mgo.Database) (*Config, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "config",
		"action":  "Init"})
	var config Config
	collection := db.C(collectionName)

	error := collection.Find(bson.M{"entity": "config"}).One(&config)

	if error != nil {
		return nil, error
	}
	config.collection = collection
	validated, validationError := govalidator.ValidateStruct(config)

	if validated != true {
		logger.Infof("Config parsing error: %+v", validationError)
		return nil, errors.New("Config parsing error")
	}
	logger.Infof("%+v", config)
	return &config, error
}
