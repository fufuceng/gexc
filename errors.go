package gexc

import "errors"

var (
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrCurrencyNotFound    = errors.New("currency not found in data")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrClientFailed        = errors.New("client failed")
)
