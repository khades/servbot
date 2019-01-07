package yandexMoney

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func GetHistory(key string, since time.Time) (*OperationHistory, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "yandexMoney",
		"action":  "operation-history"})
	apiUrl := "https://money.yandex.ru/api/operation-history"
	data := url.Values{}
	data.Set("type", "deposition")
	data.Set("details", "true")
	data.Set("from", since.Format("2006-01-02T15:04:05Z07:00"))

	client := &http.Client{Timeout: 5 * time.Second}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+key)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	dumpRequest, _ := httputil.DumpRequestOut(r, true)
	logger.Debugf("Request is: %s", dumpRequest)

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}
	dumpResponse, _ := httputil.DumpResponse(resp, true)
	logger.Debugf("Response is: %s", dumpResponse)

	result := OperationHistory{}

	marshallError := json.NewDecoder(resp.Body).Decode(&result)
	if marshallError != nil {
		return nil, marshallError
	}
	return &result, nil
}
