package models

type HttpError struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
