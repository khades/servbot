package config

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)


// Config struct describes configuration of that chatbot instance
type Config struct {
	// Dependencies
	collection               *mgo.Collection ` bson:"-"`

	// Own Fields
	APIKey           string        `valid:"required"`
	NextAPIKeyUpdate time.Time     `valid:"required"`
	OauthKey         string        `valid:"required"`
	BotUserName      string        `valid:"required"`
	BotUserID        string        `valid:"required"`
	Channels         []string
	ChannelIDs       []string
	ClientID         string `valid:"required"`
	ClientSecret     string `valid:"required"`
	AppOauthURL      string `valid:"required"`
	AppURL           string `valid:"required"`
	Debug            bool
	VkClientKey      string
	YoutubeKey       string
}


func (config *Config) SaveApiKey(key *string) {
	config.APIKey = *key
	config.collection.Update(bson.M{"entity": "config"}, bson.M{"$set": bson.M{"apikey": *key}})
}

// SaveConfigToDatabase saves current config object to database
func (config *Config) SaveConfigToDatabase() {
	config.collection.Upsert(bson.M{"entity": "config"}, bson.M{"$set": config})
}

