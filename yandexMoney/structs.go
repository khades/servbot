package yandexMoney

import "time"

type OperationHistory struct {
	Operations []Operation `json:"operations"`
}

type Operation struct {
	OperationID string    `json:"operation_id"`
	Status      string    `json:"status"`
	PatternID   string    `json:"pattern_id"`
	Direction   string    `json:"direction"`
	Amount      float64   `json:"amount"`
	DateTime    time.Time `json:"dateTime"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Comment     string    `json:"comment"`
	Details     string    `json:"details"`
	Type        string    `json:"type"`
}
