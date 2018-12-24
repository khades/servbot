package subday

import (
	"github.com/globalsign/mgo"
	"math/rand"
	"time"

	"github.com/khades/servbot/channelInfo"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/l10n"
	"github.com/sirupsen/logrus"
)

type Service struct {
	collection         *mgo.Collection
	channelInfoService *channelInfo.Service
}

// GetActive returns last active subday for specified channel
func (service *Service) GetActive(channelID *string) (*Subday, error) {
	var result Subday
	error := service.collection.Find(bson.M{
		"channelid": *channelID,
		"isactive":  true}).One(&result)
	return &result, error
}

// GetLast returns any last subday for specified channel
func (service *Service) GetLast(channelID *string) (*SubdayNoWinners, error) {
	var result SubdayNoWinners
	error := service.collection.Find(bson.M{
		"channelid": *channelID,
	}).Sort("-date").One(&result)
	return &result, error
}

// GetLastForMod returns any last subday for specified channel with extended information about it
func (service *Service) GetLastForMod(channelID *string) (*Subday, error) {
	var result Subday
	error := service.collection.Find(bson.M{
		"channelid": *channelID,
	}).Sort("-date").One(&result)
	return &result, error
}

// Search returns specified subday
func (service *Service) Get(id *string) (*SubdayNoWinners, error) {
	var result SubdayNoWinners
	error := service.collection.Find(bson.M{
		"_id": bson.ObjectIdHex(*id)}).One(&result)
	return &result, error
}

// GetForMod returns specified subday with extended information about it
func (service *Service) GetForMod(id *string) (*Subday, error) {
	var result Subday
	error := service.collection.Find(bson.M{
		"_id": bson.ObjectIdHex(*id)}).One(&result)
	return &result, error
}

// List return all subdays for specified channel
func (service *Service) List(channelID *string) ([]SubdayList, error) {
	var result []SubdayList
	error := service.collection.Find(bson.M{
		"channelid": *channelID,
	}).Sort("-date").All(&result)
	return result, error
}

// CloseAnyActive closes any active subday on specified channel
func (service *Service) CloseAnyActive(channelID *string, user *string, userID *string) {
	service.collection.UpdateAll(bson.M{
		"channelid": *channelID,
		"isactive":  true},
		bson.M{"$set": bson.M{"isactive": false, "closer": *user, "closerid": *userID}})
	service.channelInfoService.SetSubdayIsActive(channelID, false)

}

// Close closes specified subday on specified channel
func (service *Service) Close(channelID *string, id *string, user *string, userID *string) {
	service.collection.UpdateAll(bson.M{
		"channelid": *channelID,
		"_id":       bson.ObjectIdHex(*id)},
		bson.M{"$set": bson.M{"isactive": false, "closer": *user, "closerid": *userID}})
	service.channelInfoService.SetSubdayIsActive(channelID, false)

}

// PullWinner pulls specified user from winners in specified subday
func (service *Service) PullWinner(id *string, user *string) {
	service.collection.UpdateAll(bson.M{
		"_id":      bson.ObjectIdHex(*id),
		"isactive": true},
		bson.M{
			"$pull": bson.M{
				"winners": bson.M{"user": *user}}})
}

// PushWinners pushes (replaces) winners list for specified subday
func (service *Service) PushWinners(id bson.ObjectId, winners []SubdayRecord) {
	service.collection.UpdateAll(bson.M{
		"_id":      id,
		"isactive": true}, bson.M{
		"$set": bson.M{"winners": winners},
		"$push": bson.M{"winnershistory": bson.M{
			"$each": []SubdayWinnersHistory{SubdayWinnersHistory{
				Winners: winners,
				Date:    time.Now()}},
			"$sort":  bson.M{"date": -1},
			"$slice": 5}}})
}

// PickRandomWinner rolls one winner and automatically upserts him to winners in specified subday for specified channel (this was made to enforce that only streamer could roll)
func (service *Service) PickRandomWinner(channelID *string, id *string, subsOnly bool, nonsubsOnly bool) *SubdayRecord {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "subday",
		"action":  "PickRandomWinner"})
	var result Subday
	error := service.collection.Find(bson.M{"_id": bson.ObjectIdHex(*id), "channelid": *channelID, "isactive": true}).One(&result)
	if error != nil {
		logger.Debugf("some error with collection")
		return nil
	}
	votes := []SubdayRecord{}
	if len(result.Votes) == 0 {
		logger.Debugf("No Votes")

		return nil
	}
	if len(result.Votes) == len(result.Winners) {
		logger.Debugf("Nobody To Pick")

		return nil
	}
	winners := result.Winners
	for _, vote := range result.Votes {
		found := false
		for _, winner := range result.Winners {
			if vote.UserID == winner.UserID {
				found = true
				break
			}
		}
		if found == true {
			continue
		}
		if (nonsubsOnly == false && subsOnly == false) || (subsOnly == true && vote.IsSub == true) || (nonsubsOnly == true && vote.IsSub == false) {
			votes = append(votes, vote)
		}
	}
	var winner SubdayRecord
	if len(votes) == 0 {
		return nil
	}
	random := rand.Intn(len(votes))
	winner = votes[random]
	winners = append(winners, winner)
	service.PushWinners(bson.ObjectIdHex(*id), winners)
	return &winner
}

// Vote upserts specified vote variant of specified user for specified subday
func (service *Service) Vote(user *string, userID *string, isSub bool, id *bson.ObjectId, game *string) {

	service.collection.UpdateAll(bson.M{
		"_id":          *id,
		"votes.userid": userID},
		bson.M{"$set": bson.M{
			"votes.$.game":  *game,
			"votes.$.isSub": isSub,
			"votes.$.user":  *user}})

	service.collection.UpdateAll(bson.M{
		"_id":          *id,
		"votes.userid": bson.M{"$ne": userID}},
		bson.M{"$push": bson.M{"votes": SubdayRecord{
			User:   *user,
			IsSub:  isSub,
			UserID: *userID,
			Game:   *game}}})

}

// Create creates new subday IF there's no active subdays
func (service *Service) Create(channelID *string, subsOnly bool, name *string) (bool, *bson.ObjectId) {
	_, subdayError := service.GetActive(channelID)

	if subdayError == nil {
		return false, nil
	}
	id := bson.NewObjectId()
	subdayName := *name
	if subdayName == "" {
		channelInfo, channelInfoError := service.channelInfoService.GetChannelInfo(channelID)
		lang := "en"
		if channelInfoError == nil {
			lang = channelInfo.GetChannelLang()
		}
		subdayName = l10n.GetL10n(lang).SubdayCreationPrefix + time.Now().Format(time.UnixDate)
	}
	object := Subday{ID: id, Name: subdayName, SubsOnly: subsOnly, IsActive: true, Date: time.Now(), ChannelID: *channelID}
	err := service.collection.Insert(object)
	if err != nil {
		return false, nil
	}
	service.channelInfoService.SetSubdayIsActive(channelID, true)

	return true, &id
}
