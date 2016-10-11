package models

// SubAlert describes state of command Banme on chat
type SubAlert struct {
	Enabled       bool
	FirstMessage  string `bson:"firstMessage"`
	RepeatPrefix  string `bson:"repeatPrefix"`
	RepeatBody    string `bson:"repeatBody"`
	RepeatPostfix string `bson:"repeatPostfix"`
}
