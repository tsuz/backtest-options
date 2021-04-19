package model

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

// OptType is a type of option. This is either call or put
type OptType string

const (
	// DateLayout is a default date layout
	DateLayout = "2006-01-02"
)

const (
	// Call represents a call contract
	Call OptType = "call"
	// Put represents a put contract
	Put OptType = "put"
)

// OHLCV represents open, high, low, close volume for an option contract
type OHLCV struct {
	// Ask is the last contract's lowest sell price
	Ask decimal.Decimal
	// Bid is the last contract's highest buy price
	Bid decimal.Decimal
	// Close is the option contract's close price
	Close decimal.Decimal
	// Expiration is the option expiration price
	Expiration time.Time
	// High is the option contract's low price
	High decimal.Decimal
	// Low is the option contract's low price
	Low decimal.Decimal
	// Open is the option contract's open price
	Open decimal.Decimal
	// QuoteDate is the date in which this price was quoted
	QuoteDate time.Time
	// Strike is the strike price of the contract
	Strike decimal.Decimal
	// Type is the type of option contract. The value is either call or put
	Type OptType
	// UndAsk is underlying ask price
	UndAsk decimal.Decimal
	// UndBid is an underlying bid price
	UndBid decimal.Decimal
	// UndSym is an underlying symbol
	UndSym string
	// Volume is the volume traded of this contract within a specified timeframe; it's usually one day
	Volume decimal.Decimal
}

// NewOHLCV Creates a new OHLCV
func NewOHLCV(d time.Time, sym string, exp time.Time, s string, typ OptType, o, h, l, c, v, unda, undb string) (OHLCV, error) {

	close, err := decimal.NewFromString(c)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing close: %+v", c)
	}

	high, err := decimal.NewFromString(h)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing high: %+v", h)
	}

	low, err := decimal.NewFromString(l)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing low: %+v", l)
	}

	open, err := decimal.NewFromString(o)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing open: %+v", o)
	}

	strike, err := decimal.NewFromString(s)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing strike: %+v", s)
	}

	undask, err := decimal.NewFromString(unda)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing underlying ask price: %+v", unda)
	}

	undbid, err := decimal.NewFromString(undb)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing underlying big price: %+v", undb)
	}

	volume, err := decimal.NewFromString(v)
	if err != nil {
		return OHLCV{}, errors.Wrapf(err, "Error parsing underlying volume: %+v", v)
	}

	return OHLCV{
		Close:      close,
		Expiration: exp,
		High:       high,
		Low:        low,
		Open:       open,
		QuoteDate:  d,
		Strike:     strike,
		Type:       typ,
		UndAsk:     undask,
		UndBid:     undbid,
		UndSym:     sym,
		Volume:     volume,
	}, nil
}
