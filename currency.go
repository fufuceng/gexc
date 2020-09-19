package gexc

import "strings"

type Currency struct {
	Code string
	Name string
}

var supportedCurrencies = []Currency{
	{Code: "USD", Name: "us dollar"},
	{Code: "GBP", Name: "pound sterling"},
	{Code: "EUR", Name: "euro"},
	{Code: "JPY", Name: "yen"},
	{Code: "BGN", Name: "bulgarian lev"},
	{Code: "CZK", Name: "czech koruna"},
	{Code: "DKK", Name: "danish krone"},
	{Code: "HUF", Name: "forint"},
	{Code: "PLN", Name: "zloty"},
	{Code: "RON", Name: "romanian leu"},
	{Code: "SEK", Name: "swedish krona"},
	{Code: "CHF", Name: "swiss franc"},
	{Code: "ISK", Name: "iceland krona"},
	{Code: "NOK", Name: "norwegian krone"},
	{Code: "HRK", Name: "kuna"},
	{Code: "RUB", Name: "russian ruble"},
	{Code: "TRY", Name: "turkish lira"},
	{Code: "AUD", Name: "australian dollar"},
	{Code: "BRL", Name: "brazilian real"},
	{Code: "CAD", Name: "canadian dollar"},
	{Code: "CNY", Name: "yuan renminbi"},
	{Code: "HKD", Name: "hong kong dollar"},
	{Code: "IDR", Name: "rupiah"},
	{Code: "ILS", Name: "new israeli sheqel"},
	{Code: "INR", Name: "indian rupee"},
	{Code: "KRW", Name: "won"},
	{Code: "MXN", Name: "mexican peso"},
	{Code: "MYR", Name: "malaysian ringgit"},
	{Code: "NZD", Name: "new zealand dollar"},
	{Code: "PHP", Name: "philippine peso"},
	{Code: "SGD", Name: "singapore dollar"},
	{Code: "THB", Name: "baht"},
	{Code: "ZAR", Name: "rand"},
}

var (
	codeCurrencyMap = make(map[string]*Currency)
	nameCurrencyMap = make(map[string]*Currency)
)

func init() {
	for _, currency := range supportedCurrencies {
		cpy := currency

		codeCurrencyMap[currency.Code] = &cpy
		nameCurrencyMap[currency.Name] = &cpy
	}
}

func sanitizeCurrencyCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func sanitizeCurrencyName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

//CurrencyByCode converts given code to Currency object.
//code format: ISO 4217
func CurrencyByCode(code string) (Currency, bool) {
	return currencyByX(&codeCurrencyMap, sanitizeCurrencyCode(code))
}

//CurrencyByName converts given name to Currency object.
//example: euro -> EUR, forint -> HUF
func CurrencyByName(name string) (Currency, bool) {
	return currencyByX(&nameCurrencyMap, sanitizeCurrencyName(name))
}

func currencyByX(source *map[string]*Currency, key string) (Currency, bool) {
	if currency, ok := (*source)[key]; ok {
		return *currency, true
	}

	return Currency{}, false
}
