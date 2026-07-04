package tiger

import (
	"sort"

	"github.com/tigerfintech/openapi-go-sdk/model"
)

var validMarkets = []string{
	string(model.MarketAll),
	string(model.MarketUS),
	string(model.MarketHK),
	string(model.MarketCN),
	string(model.MarketSG),
}

var validBarPeriods = []string{
	string(model.BarPeriodDay),
	string(model.BarPeriodWeek),
	string(model.BarPeriodMonth),
	string(model.BarPeriodYear),
	string(model.BarPeriod1Min),
	string(model.BarPeriod5Min),
	string(model.BarPeriod15Min),
	string(model.BarPeriod30Min),
	string(model.BarPeriod60Min),
}

// IsValidMarket reports whether s is one of Tiger's known market codes.
func IsValidMarket(s string) bool {
	return contains(validMarkets, s)
}

// IsValidBarPeriod reports whether s is one of Tiger's known K-line periods.
func IsValidBarPeriod(s string) bool {
	return contains(validBarPeriods, s)
}

// ValidMarkets returns the known market codes, sorted, for use in error messages.
func ValidMarkets() []string {
	return sortedCopy(validMarkets)
}

// ValidBarPeriods returns the known K-line periods, sorted, for use in error messages.
func ValidBarPeriods() []string {
	return sortedCopy(validBarPeriods)
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func sortedCopy(list []string) []string {
	out := make([]string, len(list))
	copy(out, list)
	sort.Strings(out)
	return out
}
