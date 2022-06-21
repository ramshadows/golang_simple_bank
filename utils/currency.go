package utils

// constants for all supported currencies
const (
	KES = "KES"
	USD = "USD"
	EUR = "EUR"
)

// IsCurrencySupported returns true if the currency is supported
func IsCurrencySupported(currency string) bool {
	switch currency {
	case KES, USD, EUR:
		return true
	}
	return false
}
