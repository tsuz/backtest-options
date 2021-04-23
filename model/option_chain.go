package model

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// OptChainList is a list of option chain and quote dates
type OptChainList struct {
	// chain is a list of option chain
	quoteMap map[time.Time]*OptChain
	// quoteDates is a list of quote dates in ascending order of time
	quotes []time.Time
}

// OptChain is an option chain for specific quote date
type OptChain struct {
	QuoteDate time.Time
	UndPx     decimal.Decimal
	expiryMap map[time.Time]*OptChainExp
	expiry    []time.Time
}

// OptChainExp is an option chain for specific expiration
type OptChainExp struct {
	ExpireDate time.Time
	strike     []decimal.Decimal
	strikeMap  map[string]*OptChainStrike
}

// OptChainStrike is a row of option chain for strike, put, and call
type OptChainStrike struct {
	S    decimal.Decimal
	Put  OHLCV
	Call OHLCV
}

// NewOptionChain converts OHLCV data into option chain, OptChain
func NewOptionChain(data []OHLCV) (*OptChainList, error) {
	optChainList := &OptChainList{}
	ocs := make(map[time.Time]*OptChain)
	quoteDateList := make([]time.Time, 0)

	quoteDate := make(map[time.Time][]OHLCV)
	// extract quote date list and index OHLCV for each date
	for _, v := range data {
		if _, ok := quoteDate[v.QuoteDate]; !ok {
			quoteDate[v.QuoteDate] = make([]OHLCV, 0)
			quoteDateList = append(quoteDateList, v.QuoteDate)
		}
		quoteDate[v.QuoteDate] = append(quoteDate[v.QuoteDate], v)
	}

	optChainList.quotes = quoteDateList
	sort.Slice(optChainList.quotes, func(i, j int) bool {
		return optChainList.quotes[i].Before(optChainList.quotes[j])
	})

	// for each quote date list, retrieve expire dates
	for _, d := range quoteDateList {
		qd, ok := quoteDate[d]
		if !ok {
			return nil, errors.Errorf("Invalid quote date %+v", d)
		}
		optChain := OptChain{}

		// expiry date to ohlcv
		expiryMap := make(map[time.Time][]OHLCV)
		expiry := make([]time.Time, 0)

		var underlying *OHLCV
		for _, opt := range qd {
			// get underlying avg price
			if underlying == nil {
				if !opt.UndBid.IsZero() && !opt.UndAsk.IsZero() {
					underlying = &opt
				}
			}
			if _, ok := expiryMap[opt.Expiration]; !ok {
				expiryMap[opt.Expiration] = make([]OHLCV, 0)
				expiry = append(expiry, opt.Expiration)
			}
			expiryMap[opt.Expiration] = append(expiryMap[opt.Expiration], opt)
		}
		optChain.QuoteDate = d
		optChain.expiry = expiry
		sort.Slice(optChain.expiry, func(i, j int) bool {
			return optChain.expiry[i].Before(optChain.expiry[j])
		})
		if underlying != nil {
			optChain.UndPx = underlying.UndBid.Add(underlying.UndAsk).Div(decimal.NewFromFloat(2.0))
		}

		optExpiryMap := make(map[time.Time]*OptChainExp)
		for _, exp := range expiry {
			strikes := make([]decimal.Decimal, 0)
			strikeMap := make(map[string]*OptChainStrike)
			ohlcvs := expiryMap[exp]
			for _, ohlcv := range ohlcvs {
				if _, ok := strikeMap[ohlcv.Strike.String()]; !ok {
					strikeMap[ohlcv.Strike.String()] = &OptChainStrike{
						S: ohlcv.Strike,
					}
					strikes = append(strikes, ohlcv.Strike)
				}
				if ohlcv.Type == Call {
					strikeMap[ohlcv.Strike.String()].Call = ohlcv
				} else if ohlcv.Type == Put {
					strikeMap[ohlcv.Strike.String()].Put = ohlcv
				}
			}
			optExpiryMap[exp] = &OptChainExp{
				ExpireDate: exp,
				strike:     strikes,
				strikeMap:  strikeMap,
			}
		}
		optChain.expiryMap = optExpiryMap
		ocs[d] = &optChain
	}

	optChainList.quoteMap = ocs

	return optChainList, nil
}

// GetOptionChainForQuoteDate returns an option chain for the quote date. If strict is false, it will find a nearest date after the specified date.
func (o *OptChainList) GetOptionChainForQuoteDate(t time.Time, strict bool) *OptChain {
	if strict {
		v, _ := o.quoteMap[t]
		return v
	}
	newt := o.searchQuote(t)
	return o.quoteMap[newt]
}

// searchQuote will find a quote that is closest to this time. TODO convert this into binary search for faster lookup
func (o *OptChainList) searchQuote(t time.Time) time.Time {
	for i := 0; i < len(o.quotes); i++ {
		if o.quotes[i].Equal(t) {
			return t
		}
		if o.quotes[i].After(t) {
			return o.quotes[i]
		}
	}
	return time.Time{}
}

// GetOptionChainForExpiryDate gets expiry option chain for specific date. If strict is false, it will find a nearest date after the specified date.
func (o *OptChain) GetOptionChainForExpiryDate(t time.Time, strict bool) *OptChainExp {
	if strict {
		v, _ := o.expiryMap[t]
		return v
	}
	newt := o.searchExpiry(t)
	return o.expiryMap[newt]
}

// searchExpiry linearly searches for expiry time from time. this should be converted to binary search for faster lookup
func (o *OptChain) searchExpiry(t time.Time) time.Time {
	for i := 0; i < len(o.expiry); i++ {
		if o.expiry[i].Equal(t) {
			return t
		}
		if o.expiry[i].After(t) {
			return o.expiry[i]
		}
	}
	return time.Time{}
}

// GetOptionChainForStrike gets nearest option chain for strike. This can be below or above the value specified in the first argument.
func (o *OptChainExp) GetOptionChainForStrike(value decimal.Decimal, strict bool) *OptChainStrike {
	if strict {
		v, _ := o.strikeMap[value.String()]
		return v
	}
	news, ok := o.searchNearestStrike(value)
	if !ok {
		return nil
	}
	return o.strikeMap[news.String()]
}

// searchNearestStrike finds nearest strike value for the input value
func (o *OptChainExp) searchNearestStrike(value decimal.Decimal) (decimal.Decimal, bool) {
	if len(o.strike) == 0 {
		return value, false
	}
	for i := 0; i < len(o.strike); i++ {
		if o.strike[i].Equal(value) {
			return value, true
		}
		if o.strike[i].GreaterThan(value) {
			if i > 0 {
				if value.Sub(o.strike[i-1]).GreaterThan(
					o.strike[i].Sub(value),
				) {
					return o.strike[i], true
				} else {
					return o.strike[i-1], true
				}
			}
			return o.strike[i], true
		}
	}
	return value, false
}
