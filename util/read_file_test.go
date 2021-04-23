package util

import (
	"backtest-options/model"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestReadFile(t *testing.T) {
	jan10, _ := time.Parse(model.DateLayout, "2005-01-10")
	jan22, _ := time.Parse(model.DateLayout, "2005-01-22")

	ohlcv1, _ := model.NewOHLCV(jan10, "SPY", jan22, "130", model.Call, "1", "2", "0", "1", "1200", "118.94", "118.95")

	expData := []model.OHLCV{
		ohlcv1,
	}
	reader := NewFileReader()
	s := `underlying_symbol,quote_date,expiration,strike,option_type,open,high,low,close,trade_volume,bid_size_1545,bid_1545,ask_size_1545,ask_1545,underlying_bid_1545,underlying_ask_1545,vwap,open_interest,delivery_code
SPY,2005-01-10,2005-01-22,130.000,C,1.0000,2.0000,0.0000,1.0000,1200,0,0.0000,195,0.0500,118.9400,118.9500,0.0000,0,`

	r := csv.NewReader(strings.NewReader(s))

	data, err := reader.ReadNormalizedCSVFile(r)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error reading file"))
	}

	for idx, d := range data {
		exp := expData[idx]
		if !d.Open.Equal(exp.Open) {
			t.Errorf("Expected open to be %+v but got %+v on idx %+v", exp.Open, d.Open, idx)
		}
		if !d.Close.Equal(exp.Close) {
			t.Errorf("Expected close to be %+v but got %+v on idx %+v", exp.Close, d.Close, idx)
		}
		if !d.High.Equal(exp.High) {
			t.Errorf("Expected high to be %+v but got %+v on idx %+v", exp.High, d.High, idx)
		}
		if !d.Low.Equal(exp.Low) {
			t.Errorf("Expected low to be %+v but got %+v on idx %+v", exp.Low, d.Low, idx)
		}
		if !d.QuoteDate.Equal(exp.QuoteDate) {
			t.Errorf("Expected quote date to be %+v but got %+v on idx %+v",
				exp.QuoteDate,
				d.QuoteDate,
				idx)
		}
		if d.UndSym != exp.UndSym {
			t.Errorf("Expected underlying symbol to be %+v but got %+v on idx %+v", exp.UndSym, d.UndSym, idx)
		}
		if d.Type != exp.Type {
			t.Errorf("Expected option type to be %+v but got %+v on idx %+v", exp.Type, d.Type, idx)
		}
		if !d.Expiration.Equal(exp.Expiration) {
			t.Errorf("Expected expiration to be %+v but got %+v on idx %+v", exp.Expiration, d.Expiration, idx)
		}
		if !d.Strike.Equal(exp.Strike) {
			t.Errorf("Expected strike to be %+v but got %+v on idx %+v", exp.Strike, d.Strike, idx)
		}
		if !d.Volume.Equal(exp.Volume) {
			t.Errorf("Expected volume to be %+v but got %+v on idx %+v", exp.Volume, d.Volume, idx)
		}
		if !d.Ask.Equal(exp.Ask) {
			t.Errorf("Expected ask to be %+v but got %+v on idx %+v", exp.Ask, d.Ask, idx)
		}
		if !d.Bid.Equal(exp.Bid) {
			t.Errorf("Expected bid to be %+v but got %+v on idx %+v", exp.Bid, d.Bid, idx)
		}
	}
}
