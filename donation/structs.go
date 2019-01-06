package donation

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Donation struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	User        string        `json:"user"`
	UserID      string        `json:"userID"`
	ChannelID   string        `json:"channelID"`
	DisplayName string        `json:"displayName"`
	Message     string        `json:"message"`
	Amount      int           `json:"amount"`
	Paid        bool          `json:"paid"`
	Currency    string        `json:"currency"`
	CreatedAt   time.Time     `json:"createdAt"`
	Type        string        `json:"type"`
}
