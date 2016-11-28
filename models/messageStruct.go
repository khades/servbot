package models

import "time"

type MessageStruct struct {
	Date        time.Time
	MessageBody string
	BanLength   int
	BanReason   string
	SubCount    int
}
