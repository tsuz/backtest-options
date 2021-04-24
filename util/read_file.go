package util

import (
	"backtest-options/model"
	"encoding/csv"
	"time"

	"github.com/pkg/errors"
)

const (
	csvStdUndSym       = 0
	csvStdQuoteDate    = 1
	csvStdExp          = 2
	csvStdStrike       = 3
	csvStdOptType      = 4
	csvStdOpen         = 5
	csvStdHigh         = 6
	csvStdLow          = 7
	csvStdClose        = 8
	csvStdVol          = 9
	csvStdBidSize      = 10
	csvStdBid          = 11
	csvStdAskSize      = 12
	csvStdAsk          = 13
	csvStdUndBid       = 14
	csvStdUndAsk       = 15
	csvStdVwap         = 16
	csvStdOpenInterest = 17
	csvStdDelivCode    = 18
)

// MyReader is a reader interface
type MyReader interface {
	ReadNormalizedCSVFile(r *csv.Reader) ([]model.OHLCV, error)
}

type fr struct{}

func (fr *fr) ReadNormalizedCSVFile(r *csv.Reader) ([]model.OHLCV, error) {

	ohlcvs := make([]model.OHLCV, 0)

	fields, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "Error reading all file values")
	}
	for row, field := range fields {
		if row == 0 {
			continue
		}
		if len(field) < csvStdUndAsk+1 {
			return nil, errors.Errorf("Expected at least %+v rows but got %+v on row: %d",
				csvStdUndAsk+1,
				len(field),
				row+1)
		}

		optType := field[csvStdOptType]
		var typ model.OptType
		if optType == "C" {
			typ = model.Call
		} else if optType == "P" {
			typ = model.Put
		}
		quoteDate := field[csvStdQuoteDate]
		quoteTime, err := time.Parse(model.DateLayout, quoteDate)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing quote date %+v at row: %d", quoteDate, row+1)
		}
		expDate := field[csvStdExp]
		expTime, err := time.Parse(model.DateLayout, expDate)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing exp date %+v at row: %d", expDate, row+1)
		}
		ohlcv, err := model.NewOHLCV(
			quoteTime,
			field[csvStdUndSym],
			expTime,
			field[csvStdStrike],
			typ,
			field[csvStdOpen],
			field[csvStdHigh],
			field[csvStdLow],
			field[csvStdClose],
			field[csvStdVol],
			field[csvStdAsk],
			field[csvStdBid],
			field[csvStdUndAsk],
			field[csvStdUndBid],
		)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting into OHLCV")
		}

		ohlcvs = append(ohlcvs, ohlcv)

	}

	return ohlcvs, nil
}

// NewFileReader generates MyReader
func NewFileReader() MyReader {
	return &fr{}
}
