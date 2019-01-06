package currencyConverter

import (
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

type Service struct {
	rates []Rate
}

func (service *Service) get() error {
	var timeout = 5 * time.Second
	var httpClient = http.Client{Timeout: timeout}
	resp, err := httpClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()
	}
	result := RateResponse{}
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	marshallError := decoder.Decode(&result)
	if marshallError != nil {
		return marshallError
	}
	log.Printf("%+v", result)

	if len(result.Cube) == 0 {
		return errors.New("No response")
	}

	service.rates = result.Cube[0].Rates
	return nil
}
func (service *Service) ConvertToUSD(amount float64, currencyCode string) float64 {
	if currencyCode == "USD" {
		return amount
	}

	defaultRate := Rate{}
	foundRate := Rate{}
	for _, rate := range service.rates {
		if rate.Currency == currencyCode {
			foundRate = rate
			break
		}
		if rate.Currency == "USD" {
			defaultRate = rate
		}
	}
	if foundRate.Currency != "" {
		return amount / (foundRate.Rate) * (defaultRate.Rate)
	} else {
		return amount * (defaultRate.Rate)
	}
}
