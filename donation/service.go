package donation

import (
	"time"

	"github.com/khades/servbot/balance"
	"github.com/khades/servbot/currencyConverter"
	"github.com/khades/servbot/event"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	collection               *mgo.Collection
	currencyConverterService *currencyConverter.Service
	balanceService           *balance.Service
	eventService             *event.Service
}

func (service *Service) Create(user string, userID string, displayName string, channelID string, amount int, currency string) *bson.ObjectId {
	objectID := bson.NewObjectId()
	donation := Donation{
		ID:          objectID,
		User:        user,
		UserID:      userID,
		DisplayName: displayName,
		ChannelID:   channelID,
		Amount:      amount,
		Currency:    currency,
		CreatedAt:   time.Now(),
		Type:        "donation",
	}
	service.collection.Insert(donation)
	return &objectID
}

func (service *Service) List(channelID string, page int) ([]Donation, error) {
	results := []Donation{}
	const pageSize int = 50
	error := service.collection.Find(bson.M{"paid": true, "shown": false}).Sort("-date").Skip((page - 1) * pageSize).Limit(pageSize).All(&results)
	return results, error
}

func (service *Service) SetPaid(id string, channelID string) (string, error) {
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"paid": true}},
		ReturnNew: true,
	}
	result := Donation{}
	_, error := service.collection.Find(bson.M{"_id": bson.ObjectIdHex(id), "channelid": channelID}).Apply(change, &result)

	if error == nil {
		service.eventService.Put(result.ChannelID, event.Event{
			User:     result.User,
			Type:     event.DONATION,
			Amount:   result.Amount,
			Message:  result.Message,
			Currency: result.Currency,
		})

		service.balanceService.Inc(
			result.ChannelID,
			result.UserID,
			result.UserID,
			service.currencyConverterService.ConvertToUSD(float64(result.Amount), result.Currency),
		)
	}
	return result.ChannelID, error
}
