package openex

type Config struct {
	BaseUrl  string
	Protocol string
	Port     string
}

var defaultConfig = Config{
	BaseUrl:  "api.exchangeratesapi.io",
	Protocol: "https",
}
