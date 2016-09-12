package models

// SubAlert describes state of command Banme on chat
type SubAlert struct {
	Enabled       bool
	FirstMessage  string
	RepeatPrefix  string
	RepeatBody    string
	RepeatPostfix string
}
