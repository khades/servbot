package repos

import "gopkg.in/mgo.v2/bson"

func SetTrackPlayed(id string) {

	Db.C("songRequests").Update(
		bson.M{
			"_id": bson.ObjectId(id)},
		bson.M{"$set": bson.M{
			"inqueue":    false,
			"playingnow": false,
		}})

}
