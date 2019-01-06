package config

import (
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Config struct describes configuration of that chatbot instance
type Config struct {
	// Dependencies
	collection *mgo.Collection ` bson:"-"`

	// Own Fields
	APIKey             string
	NextAPIKeyUpdate   time.Time
	OauthKey           string   `valid:"required"`
	BotUserName        string   `valid:"required"`
	BotUserID          string   `valid:"required"`
	Channels           []string ` bson:"-"`
	ChannelIDs         []string ` bson:"-"`
	ClientID           string   `valid:"required"`
	ClientSecret       string   `valid:"required"`
	AppOauthURL        string   `valid:"required"`
	AppURL             string   `valid:"required"`
	YandexClientID     string
	YandexClientSecret string
	Debug              bool
	VkClientKey        string
	YoutubeKey         string
}

func (config *Config) SaveApiKey(key string, nextApiKeyUpdate time.Time) {
	config.APIKey = key
	config.NextAPIKeyUpdate = nextApiKeyUpdate

	config.collection.Update(bson.M{"entity": "config"},
		bson.M{"$set": bson.M{
			"apikey":           key,
			"nextapikeyupdate": nextApiKeyUpdate,
		}})
}

func (config *Config) NeedsAPIKey(period time.Duration) bool {
	return strings.TrimSpace(config.APIKey) == "" ||
		time.Now().Add(period).After(config.NextAPIKeyUpdate)
}

// SaveConfigToDatabase saves current config object to database
func (config *Config) SaveConfigToDatabase() {
	config.collection.Upsert(bson.M{"entity": "config"}, bson.M{"$set": config})
}
