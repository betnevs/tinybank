package valid

const (
	CNY = "CNY"
	EUR = "EUR"
	USD = "USD"
)

func IsSupportCurrency(currency string) bool {
	switch currency {
	case EUR, CNY, USD:
		return true
	}
	return false
}
