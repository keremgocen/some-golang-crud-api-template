package mockexchange

import "fmt"

// Converter implements the exchange rate converter behaviour.
type Converter struct {
	// any config
}

// NewConverter returns a new converter.
func NewConverter() *Converter {
	return &Converter{}
}

// ConvertExchangeRate returns the live exchange rate value between currencies.
func (c *Converter) ConvertExchangeRate(from, to string) (float64, error) {
	// lookup the rates from the imaginary data store
	rates, err := imaginaryDataStoreLookUp(from)
	if err != nil {
		return 0, fmt.Errorf("data store err:%w", err)
	}

	r, ok := rates[to]
	if !ok {
		return 0, fmt.Errorf("currency not found for %s", to)
	}

	return r, nil
}

// imaginaryDataStoreLookUp mimics a live exchange rate API lookup.
func imaginaryDataStoreLookUp(countryCode string) (map[string]float64, error) {
	if countryCode == "GBP" {
		exchangeRateMap := make(map[string]float64)
		exchangeRateMap["EUR"] = 1.19
		exchangeRateMap["USD"] = 1.21
		exchangeRateMap["CAD"] = 1.56
		exchangeRateMap["SEK"] = 12.46
		exchangeRateMap["SEK"] = 160.55
		return exchangeRateMap, nil
	}
	return nil, fmt.Errorf("countryCode:%s data not available", countryCode)
}
