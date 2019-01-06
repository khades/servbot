package yandexMoney

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetHistory(key string, since time.Time) (*OperationHistory, error) {
	apiUrl := "https://money.yandex.ru/api/operation-history"
	data := url.Values{}
	data.Set("type", "deposition")
	data.Set("details", "true")
	data.Set("from", since.Format("2006-01-02T15:04:05Z07:00"))

	client := &http.Client{timeout: 5 * time.Second}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+key)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}

	result := OperationHistory{}

	marshallError := json.NewDecoder(resp.Body).Decode(&result)
	if marshallError != nil {
		return nil, marshallError
	}
	return &result, nil
}
