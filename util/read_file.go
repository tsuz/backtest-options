package util

import (
	"archive/zip"
	"encoding/csv"
	"option-analysis/model"
	"time"

	"github.com/pkg/errors"
)

const (
	csvUndSym    = 0
	csvQuoteDate = 1
	csvExp       = 3
	csvStrike    = 4
	csvOptType   = 5
	csvOpen      = 6
	csvHigh      = 7
	csvLow       = 8
	csvClose     = 9
	csvVol       = 10
	csvBid       = 18
	csvAsk       = 20
	csvUndBid    = 21
	csvUndAsk    = 22
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
			if len(field) < csvUndAsk+1 {
				return nil, errors.Errorf("Expected at least %+v rows but got %+v on row: %d", csvUndAsk+1, len(field), row+1)
			}

			optType := field[csvOptType]
			var typ model.OptType
			if optType == "C" {
				typ = model.Call
			} else if optType == "P" {
				typ = model.Put
			}
			quoteDate := field[csvQuoteDate]
			quoteTime, err := time.Parse(model.DateLayout, quoteDate)
			if err != nil {
				return nil, errors.Wrapf(err, "Error parsing quote date %+v at row: %d", quoteDate, row+1)
			}
			expDate := field[csvExp]
			expTime, err := time.Parse(model.DateLayout, expDate)
			if err != nil {
				return nil, errors.Wrapf(err, "Error parsing exp date %+v at row: %d", expDate, row+1)
			}
			ohlcv, err := model.NewOHLCV(
				quoteTime,
				field[csvUndSym],
				expTime,
				field[csvStrike],
				typ,
				field[csvOpen],
				field[csvHigh],
				field[csvLow],
				field[csvClose],
				field[csvVol],
				field[csvUndAsk],
				field[csvUndBid],
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
