package util

import (
	"option-analysis/model"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestReadFile(t *testing.T) {
	june1, _ := time.Parse(model.DateLayout, "2016/06/01")
	june8, _ := time.Parse(model.DateLayout, "2016/06/08")

	ohlcv1, _ := model.NewOHLCV(june1, "^VIX", june8, "14", model.Put, "0.1", "0.25", "0.1", "0.25", "623", "14.24", "14.24")
	ohlcv2, _ := model.NewOHLCV(june1, "^VIX", june8, "14.5", model.Call, "1", "1.1", "0.65", "0.65", "55", "14.24", "14.24")
	ohlcv3, _ := model.NewOHLCV(june1, "^VIX", june8, "14.5", model.Put, "0.3", "0.55", "0.3", "0.55", "67", "14.24", "14.24")
	ohlcv4, _ := model.NewOHLCV(june1, "^VIX", june8, "15", model.Call, "0.9", "0.95", "0.5", "0.5", "82", "14.24", "14.24")

	expData := []model.OHLCV{
		ohlcv1,
		ohlcv2,
		ohlcv3,
		ohlcv4,
	}
	reader := NewFileReader()
	data, err := reader.ReadFile("./../testdata/quotes_sample.zip")
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error reading file"))
	}
	if len(data) != 4 {
		t.Errorf("Expected 4 items but got %+v", len(data))
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
