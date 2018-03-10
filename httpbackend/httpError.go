package httpbackend

type httpError struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
