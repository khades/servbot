package currencyConverter

func Init() *Service {
	service := Service{}
	service.get()
	return &service
}
