package util

import (
	"archive/zip"
	"encoding/csv"
	"io/ioutil"

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

type liveVolImporter struct {
	writer *csv.Writer
}

// NewLiveVolImporter is a live vol importer
func NewLiveVolImporter(w *csv.Writer) Importer {
	return &liveVolImporter{
		writer: w,
	}
}

// ImportFolder imports a folder and outputs to specified data directory
func (livevol *liveVolImporter) ImportFolder(folder string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return errors.Wrap(err, "Error importing folder")
	}

	header := []string{
		"underlying_symbol",
		"quote_date",
		"expiration",
		"strike",
		"option_type",
		"open",
		"high",
		"low",
		"close",
		"trade_volume",
		"bid_size_1545",
		"bid_1545",
		"ask_size_1545",
		"ask_1545",
		"underlying_bid_1545",
		"underlying_ask_1545",
		"vwap",
		"open_interest",
		"delivery_code",
	}
	livevol.writer.Write(header)

	for _, file := range files {
		f, err := zip.OpenReader(folder + "/" + file.Name())
		if err != nil {
			return errors.Wrapf(err, "Error opening file %s", folder+"/"+file.Name())
		}
		defer f.Close()

		for _, file := range f.File {
			fopen, err := file.Open()
			if err != nil {
				return errors.Wrapf(err, "Error opening file %+v", file.Name)
			}

			reader := csv.NewReader(fopen)
			fields, err := reader.ReadAll()
			if err != nil {
				return errors.Wrapf(err, "Error reading all %+v", file.Name)
			}

			for row, field := range fields {
				if row == 0 {
					continue
				}
				if len(field) < csvLivevolDelivCode {
					continue
				}

				values := []string{
					field[csvLivevolUndSym],
					field[csvLivevolQuoteDate],
					field[csvLivevolExp],
					field[csvLivevolStrike],
					field[csvLivevolOptType],
					field[csvLivevolOpen],
					field[csvLivevolHigh],
					field[csvLivevolLow],
					field[csvLivevolClose],
					field[csvLivevolVol],
					field[csvLivevolBidSize],
					field[csvLivevolBid],
					field[csvLivevolAskSize],
					field[csvLivevolAsk],
					field[csvLivevolUndBid],
					field[csvLivevolUndAsk],
					field[csvLivevolVwap],
					field[csvLivevolOpenInterest],
					field[csvLivevolDelivCode],
				}
				livevol.writer.Write(values)
			}
		}
	}
	return nil
}

// ImportFile imports a file and outputs to specified data directory
func (livevol *liveVolImporter) ImportFile(file string, output string) error {
	return nil
}
