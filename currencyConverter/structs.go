package currencyConverter

type Rate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

type RateResponse struct {
	Cube []struct {
		Date  string `xml:"time,attr"`
		Rates []Rate `xml:"Cube"`
	} `xml:"Cube>Cube"`
}
