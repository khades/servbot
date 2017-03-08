package models

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
