package mockexchange

// CurrencyConverter is a mock interface to provide exchange rate API.
type CurrencyConverter interface {
	// ConvertExchangeRate returns the exchange rate between provided currencies
	ConvertExchangeRate(from, to string) (float64, error)
}
