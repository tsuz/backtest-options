package model

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

func TestOptionChain(t *testing.T) {
	may1, _ := time.Parse(DateLayout, "2016/05/01")
	june1, _ := time.Parse(DateLayout, "2016/06/01")
	june20, _ := time.Parse(DateLayout, "2016/06/20")
	july2, _ := time.Parse(DateLayout, "2016/07/02")
	aug1, _ := time.Parse(DateLayout, "2016/08/01")
	ohlcv1, _ := NewOHLCV(june1, "SPY", july2, "116", Call, "1", "1", "1", "1", "623", "115.5", "116.5")
	ohlcv2, _ := NewOHLCV(june1, "SPY", july2, "116", Put, "0.5", "0.5", "0.5", "0.5", "623", "115.5", "116.5")
	ohlcv3, _ := NewOHLCV(july2, "SPY", july2, "116", Call, "0.5", "0.5", "0.5", "0.5", "623", "117.5", "118.5")
	testData := []OHLCV{
		ohlcv1,
		ohlcv2,
		ohlcv3,
	}

	chain, err := NewOptionChain(testData)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating new option chain"))
	}

	// should not get option chain for quote date if there are no quotes on or after this date
	augchain := chain.GetOptionChainForQuoteDate(aug1, false)
	if augchain != nil {
		t.Error(errors.Errorf("Expected to be nil but got %+v", augchain))
	}

	// should not get optoin chain for quote date if it's strict and it doesn't exist on that date
	maychain := chain.GetOptionChainForQuoteDate(may1, true)
	if maychain != nil {
		t.Error(errors.Errorf("Expected to be nil but got %+v", maychain))
	}

	// should get option chain for quote date since data exists
	oc := chain.GetOptionChainForQuoteDate(june1, false)
	if oc == nil {
		t.Error(errors.Errorf("Error getting option chain for %+v", june1))
	}

	if !oc.QuoteDate.Equal(june1) {
		t.Error(errors.Errorf("Expected quote time to be %+v but got %+v", june1, oc.QuoteDate))
	}

	if !oc.UndPx.Equal(decimal.NewFromFloat(116)) {
		t.Error(errors.Errorf("Expected underlying price to be %+v but got %+v", 116, oc.UndPx.String()))
	}

	// should not get options for expiry date if it's strict and data doesn't exist
	nooc1 := oc.GetOptionChainForExpiryDate(june1, true)
	if nooc1 != nil {
		t.Error(errors.Errorf("Expected option chain to be nil for expiry %+v but got %+v", june1, nooc1))
	}

	// should not get options for expiry date if it's not strict and data doesn't exist after that date
	nooc2 := oc.GetOptionChainForExpiryDate(aug1, false)
	if nooc2 != nil {
		t.Error(errors.Errorf("Expected option chain to be nil for expiry %+v but got %+v", aug1, nooc2))
	}

	expchain1 := oc.GetOptionChainForExpiryDate(june20, false)
	expchain2 := oc.GetOptionChainForExpiryDate(july2, false)
	if expchain1 != expchain2 {
		t.Error(errors.Errorf("Expected two options to be the same since strict is false but got chain 1 %+v vs chain 2 %+v", expchain1, expchain2))
	}
	if !expchain1.ExpireDate.Equal(july2) {
		t.Error(errors.Errorf("Expected date to be %+v but got %+v", july2, expchain1.ExpireDate))
	}

	farstrike, _ := decimal.NewFromString("115")
	validstrike, _ := decimal.NewFromString("116")
	// if it's strict, this should return nothing
	nos := expchain1.GetOptionChainForStrike(farstrike, true)
	if nos != nil {
		t.Error(errors.Errorf("Expected option chain to be nil for expiry %+v and strike %+v but got %+v", june1, farstrike, nos))
	}
	valid1 := expchain1.GetOptionChainForStrike(farstrike, false)
	valid2 := expchain1.GetOptionChainForStrike(validstrike, true)
	valid3 := expchain1.GetOptionChainForStrike(validstrike, false)
	if valid1 != valid2 {
		t.Error(errors.Errorf("Expected option chain (1 and 2) to be same for strike %+v and %+v but got %+v vs %+v",
			farstrike,
			validstrike,
			valid1,
			valid2,
		))
	}
	if valid2 != valid3 {
		t.Error(errors.Errorf("Expected option chain (2 and 3) to be same for strike %+v and %+v but got %+v vs %+v",
			validstrike,
			validstrike,
			valid2,
			valid3,
		))
	}

	if valid3.Call != ohlcv1 {
		t.Error(errors.Errorf("Expected call ohlcv and initial ohlcv to be the same but got %+v vs %+v",
			valid3.Call,
			ohlcv1,
		))
	}

	if valid3.Put != ohlcv2 {
		t.Error(errors.Errorf("Expected put ohlcv and initial ohlcv to be the same but got %+v vs %+v",
			valid3.Put,
			ohlcv2,
		))
	}
}
