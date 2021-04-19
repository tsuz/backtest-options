package model

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestOHLCV(t *testing.T) {
	layout := DateLayout
	qdstr := "2005/01/13"
	quoteDate, err := time.Parse(layout, qdstr)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error quote date failed to parse"))
	}
	edstr := "2005/01/13"
	expi, err := time.Parse(layout, edstr)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error quote date failed to parse"))
	}
	symbol := "SPY"
	open := "20.6"
	high := "20.7"
	low := "20.4"
	close := "20.5"
	vol := "50"
	strike := "117.5"
	undAsk := "117.7"
	undBid := "117.6"
	ohlcv, _ := NewOHLCV(quoteDate, symbol, expi, strike, Call, open, high, low, close, vol, undAsk, undBid)
	if ohlcv.QuoteDate != quoteDate {
		t.Error(errors.Errorf("Expected %+v but got %+v", quoteDate, ohlcv.QuoteDate))
	}
	if ohlcv.UndSym != symbol {
		t.Error(errors.Errorf("Expected %+v but got %+v", symbol, ohlcv.UndSym))
	}
	if !ohlcv.Expiration.Equal(expi) {
		t.Error(errors.Errorf("Expected %+v but got %+v", expi, ohlcv.Expiration))
	}
	if !ohlcv.QuoteDate.Equal(quoteDate) {
		t.Error(errors.Errorf("Expected %+v but got %+v", quoteDate, ohlcv.QuoteDate))
	}
	if ohlcv.Open.String() != open {
		t.Error(errors.Errorf("Expected %+v but got %+v", open, ohlcv.Open))
	}
	if ohlcv.Close.String() != close {
		t.Error(errors.Errorf("Expected %+v but got %+v", close, ohlcv.Close))
	}
	if ohlcv.High.String() != high {
		t.Error(errors.Errorf("Expected %+v but got %+v", high, ohlcv.High))
	}
	if ohlcv.Low.String() != low {
		t.Error(errors.Errorf("Expected %+v but got %+v", low, ohlcv.Low))
	}
	if ohlcv.Strike.String() != strike {
		t.Error(errors.Errorf("Expected %+v but got %+v", strike, ohlcv.Strike))
	}
	if ohlcv.Volume.String() != vol {
		t.Error(errors.Errorf("Expected %+v but got %+v", vol, ohlcv.Volume))
	}
	if ohlcv.UndAsk.String() != undAsk {
		t.Error(errors.Errorf("Expected %+v but got %+v", undAsk, ohlcv.UndAsk))
	}
	if ohlcv.UndBid.String() != undBid {
		t.Error(errors.Errorf("Expected %+v but got %+v", undBid, ohlcv.UndBid))
	}
}
