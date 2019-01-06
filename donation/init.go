package donation

import (
	"github.com/globalsign/mgo"
	"github.com/khades/servbot/balance"
	"github.com/khades/servbot/currencyConverter"
	"github.com/khades/servbot/event"
)

var collectionName = "donations"

func Init(db *mgo.Database,
	currencyConverterService *currencyConverter.Service,
	balanceService *balance.Service,
	eventService *event.Service) *Service {
	collection := db.C(collectionName)

	collection.EnsureIndex(mgo.Index{
		Key: []string{"channelid", "paid", "shown"}})

	return &Service{
		collection:               collection,
		balanceService:           balanceService,
		currencyConverterService: currencyConverterService,
		eventService:             eventService,
	}
}
