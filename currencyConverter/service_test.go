package currencyConverter

import "testing"

func TestIsWorking(t *testing.T) {
	service := Service{}
	err := service.get()
	if err != nil {
		t.Logf("%+v", err)

		t.Fail()
	}
	t.Fail()
}

func TestConvertProperlyFromRouble(t *testing.T) {
	service := Service{[]Rate{
		Rate{
			Currency: "USD",
			Rate:     2,
		}, Rate{
			Currency: "RUB",
			Rate:     35,
		}}}
	if service.ConvertToUSD(140.000, "RUB") != 8.0 {
		t.Fail()
	}
}

func TestConvertProperlyFromEuro(t *testing.T) {
	service := Service{[]Rate{
		Rate{
			Currency: "USD",
			Rate:     4,
		}}}
	if service.ConvertToUSD(2.000, "EUR") != 8.0 {
		t.Fail()
	}
}
