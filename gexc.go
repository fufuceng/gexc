package gexc

import (
	"fmt"
	"github.com/fufuceng/gexc/internal/openex"
	"github.com/fufuceng/gexc/response"
	time2 "github.com/fufuceng/gexc/time"
	"time"
)

type fxToWrapper struct {
	base   *Fx
	from   string
	amount float64
}

//To is the last step of the currency conversion.
//Takes currency that base currency will be converted to it.
//Validates all parameters and converts one currency to another.
//May raises validation and connection errors.
func (f *fxToWrapper) To(currency string) (float64, error) {
	fromCurrency, ok := CurrencyByCode(f.from)
	if !ok {
		return 0, fmt.Errorf("%w: %v", ErrUnsupportedCurrency, f.from)
	}

	toCurrency, ok := CurrencyByCode(currency)
	if !ok {
		return 0, fmt.Errorf("%w: %v", ErrUnsupportedCurrency, currency)
	}

	resp, err := f.base.HistoryOf(fromCurrency.Code).Against(toCurrency.Code).Latest()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrClientFailed, err)
	}

	mul, ok := resp.Rates[toCurrency.Code]
	if !ok {
		return 0, fmt.Errorf("%w: %v", ErrCurrencyNotFound, toCurrency.Code)
	}

	return f.amount * mul, nil
}

type fxFromWrapper struct {
	base   *Fx
	amount float64
}

//From is the second step of the currency conversion
//It takes base currency that will be converted to the another.
func (f *fxFromWrapper) From(currency string) *fxToWrapper {
	return &fxToWrapper{
		base:   f.base,
		from:   currency,
		amount: f.amount,
	}
}

type fxHistoryUntilWrapper struct {
	base     *Fx
	currency string
	against  []string
	from     time.Time
}

//Until is the last step of the history
//It takes time value and validates all parameters then returns the history.
func (f *fxHistoryUntilWrapper) Until(t time.Time) (response.History, error) {
	if f.from.Equal(time.Time{}) || t.Equal(time.Time{}) {
		return response.History{}, fmt.Errorf("%w: time values should not be empty", ErrInvalidParameter)
	}

	if !t.After(f.from) {
		return response.History{}, fmt.Errorf("%w: until value should be bigger than from", ErrInvalidParameter)
	}

	currency, ok := CurrencyByCode(f.currency)
	if !ok {
		return response.History{}, fmt.Errorf("unsuportted currency: %v", f.currency)
	}

	var againstCurrencies []string
	for _, code := range f.against {
		cur, ok := CurrencyByCode(code)
		if !ok {
			return response.History{}, fmt.Errorf("%w: %v", ErrUnsupportedCurrency, code)
		}

		againstCurrencies = append(againstCurrencies, cur.Code)
	}

	resp, err := f.base.openexClient.History(openex.HistoryParams{
		StartAt: time2.NewGexc(f.from),
		EndAt:   time2.NewGexc(t),
		Base:    currency.Code,
		Symbols: againstCurrencies,
	})

	if err != nil {
		return response.History{}, err
	}

	return *resp, nil
}

type fxHistoryFromWrapper struct {
	base         *Fx
	baseCurrency string
	against      []string
}

//From is the third step of the history.
//It takes initial date
func (f *fxHistoryFromWrapper) From(time time.Time) *fxHistoryUntilWrapper {
	return &fxHistoryUntilWrapper{
		base:     f.base,
		currency: f.baseCurrency,
		against:  f.against,
		from:     time,
	}
}

//Latest calculates currency values corresponding to the given currency based on today
func (f *fxHistoryFromWrapper) Latest() (response.SingleDate, error) {
	baseCurrency, ok := CurrencyByCode(f.baseCurrency)
	if !ok {
		return response.SingleDate{}, fmt.Errorf("%w: %v", ErrUnsupportedCurrency, f.baseCurrency)
	}

	var againstCurrencies []string
	for _, code := range f.against {
		cur, ok := CurrencyByCode(code)
		if !ok {
			return response.SingleDate{}, fmt.Errorf("%w: %v", ErrUnsupportedCurrency, code)
		}

		againstCurrencies = append(againstCurrencies, cur.Code)
	}

	resp, err := f.base.openexClient.Latest(openex.LatestParams{
		Base:    baseCurrency.Code,
		Symbols: againstCurrencies,
	})

	if err != nil {
		return response.SingleDate{}, fmt.Errorf("%w: %v", ErrClientFailed, err)
	}

	return *resp, nil
}

type fxHistoryWrapper struct {
	base         *Fx
	baseCurrency string
}

//Against is the second step of the currency history
//It takes currencies that will be compared to base history in time range
func (f *fxHistoryWrapper) Against(currencies ...string) *fxHistoryFromWrapper {
	return &fxHistoryFromWrapper{
		base:         f.base,
		baseCurrency: f.baseCurrency,
		against:      currencies,
	}
}

//Fx collects all functionality of the library
//It includes Amount, Convert and HistoryOf functions
type Fx struct {
	openexClient openex.Client
}

//Amount is the initial step of the currency conversion.
//It takes amount that will be converted.
func (f *Fx) Amount(amount float64) *fxFromWrapper {
	return &fxFromWrapper{
		base:   f,
		amount: amount,
	}
}

//Convert is the short form of `Amount.From.To` chain.
//It calculates the final currency that corresponding to base currency
func (f *Fx) Convert(amount float64, from, to string) (float64, error) {
	return f.Amount(amount).From(from).To(to)
}

//HistoryOf is the initial step of the collection of the currency history
//It takes base currency that will be compared to others in time range
func (f *Fx) HistoryOf(currency string) *fxHistoryWrapper {
	return &fxHistoryWrapper{
		base:         f,
		baseCurrency: currency,
	}
}

func New() *Fx {
	return &Fx{
		openexClient: openex.NewDefaultClient(),
	}
}

func newWithClient(client openex.Client) *Fx {
	return &Fx{openexClient: client}
}
