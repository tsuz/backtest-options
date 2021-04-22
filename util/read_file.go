package util

import (
	"archive/zip"
	"backtest-options/model"
	"encoding/csv"
	"time"

	"github.com/pkg/errors"
)

const (
	csvLivevolUndSym       = 0
	csvLivevolQuoteDate    = 1
	csvLivevolExp          = 3
	csvLivevolStrike       = 4
	csvLivevolOptType      = 5
	csvLivevolOpen         = 6
	csvLivevolHigh         = 7
	csvLivevolLow          = 8
	csvLivevolClose        = 9
	csvLivevolVol          = 10
	csvLivevolBidSize      = 11
	csvLivevolBid          = 12
	csvLivevolAskSize      = 13
	csvLivevolAsk          = 14
	csvLivevolUndBid       = 15
	csvLivevolUndAsk       = 16
	csvLivevolVwap         = 23
	csvLivevolOpenInterest = 24
	csvLivevolDelivCode    = 25
)

// MyReader is a reader interface
type MyReader interface {
	ReadFile(file string) ([]model.OHLCV, error)
}

type fr struct{}

func (r *fr) ReadFile(file string) ([]model.OHLCV, error) {

	f, err := zip.OpenReader(file)
	if err != nil {
		return nil, errors.Wrapf(err, "Error opening file %s", file)
	}

	defer f.Close()

	ohlcvs := make([]model.OHLCV, 0)

	for _, file := range f.File {
		fopen, err := file.Open()
		if err != nil {
			return nil, errors.Wrapf(err, "Error opening file %+v", file.Name)
		}

		reader := csv.NewReader(fopen)
		fields, err := reader.ReadAll()
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading all %+v", file.Name)
		}
		for row, field := range fields {
			if row == 0 {
				continue
			}
			if len(field) < csvLivevolUndAsk+1 {
				return nil, errors.Errorf("Expected at least %+v rows but got %+v on row: %d",
					csvLivevolUndAsk+1,
					len(field),
					row+1)
			}

			optType := field[csvLivevolOptType]
			var typ model.OptType
			if optType == "C" {
				typ = model.Call
			} else if optType == "P" {
				typ = model.Put
			}
			quoteDate := field[csvLivevolQuoteDate]
			quoteTime, err := time.Parse(model.DateLayout, quoteDate)
			if err != nil {
				return nil, errors.Wrapf(err, "Error parsing quote date %+v at row: %d", quoteDate, row+1)
			}
			expDate := field[csvLivevolExp]
			expTime, err := time.Parse(model.DateLayout, expDate)
			if err != nil {
				return nil, errors.Wrapf(err, "Error parsing exp date %+v at row: %d", expDate, row+1)
			}
			ohlcv, err := model.NewOHLCV(
				quoteTime,
				field[csvLivevolUndSym],
				expTime,
				field[csvLivevolStrike],
				typ,
				field[csvLivevolOpen],
				field[csvLivevolHigh],
				field[csvLivevolLow],
				field[csvLivevolClose],
				field[csvLivevolVol],
				field[csvLivevolUndAsk],
				field[csvLivevolUndBid],
			)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting into OHLCV")
			}

			ohlcvs = append(ohlcvs, ohlcv)
		}
	}

	return ohlcvs, nil
}

// NewFileReader generates MyReader
func NewFileReader() MyReader {
	return &fr{}
}
