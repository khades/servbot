package donationSource

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	collection *mgo.Collection
}

func (service *Service) SetYandexKey(channelID string, key string) {
	// TODO Locate expiration date
	service.collection.Upsert(bson.M{"channelid": channelID}, bson.M{"$set": bson.M{"yandex": DonationSource{
		Enabled:        true,
		Key:            key,
		ExpirationDate: time.Now().Add(3 * 365 * 24 * time.Hour),
		LastCheck:      time.Now(),
	}}})
}

func (service *Service) Get(channelID string) (*DonationSources, error) {
	result := DonationSources{}
	error := service.collection.Find(bson.M{"channelid": channelID}).One(&result)
	return &result, error
}

func (service *Service) List() ([]DonationSources, error) {
	result := []DonationSources{}
	error := service.collection.Find(bson.M{}).All(&result)
	return result, error
}

func (service *Service) UpdateYandexLastCheck(channelID string, lastcheck time.Time) {
	service.collection.Update(bson.M{"channelid": channelID}, bson.M{"$set": bson.M{"yandex.lastcheck": lastcheck}})
}
