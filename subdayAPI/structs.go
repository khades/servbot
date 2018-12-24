package subdayAPI

import "github.com/khades/servbot/subday"

type subdayWithMod struct {
	*subday.Subday
	IsMod bool `json:"isMod"`
}

type subdayWithModNoWinners struct {
	*subday.SubdayNoWinners
	IsMod bool `json:"isMod"`
}

type subdayCreateStruct struct {
	Name     string `json:"name"`
	SubsOnly bool   `json:"subsOnly"`
}

type subdayCreateResp struct {
	ID string `json:"id"`
}