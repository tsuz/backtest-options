package util

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestImportFile(t *testing.T) {

	var b bytes.Buffer
	w := csv.NewWriter(&b)
	w.UseCRLF = false

	imp := NewLiveVolImporter(w)
	if err := imp.ImportFolder("./../testdata"); err != nil {
		t.Error(errors.Wrap(err, "Error from import folder func"))
	}
	w.Flush()
	if err := w.Error(); err != nil {
		t.Fatal(errors.Wrap(err, "Error from csv flush"))
	}

	expectedRows := []string{
		"underlying_symbol,quote_date,expiration,strike,option_type,open,high,low,close,trade_volume,bid_size_1545,bid_1545,ask_size_1545,ask_1545,underlying_bid_1545,underlying_ask_1545,vwap,open_interest,delivery_code",
		"^VIX,2016-06-03,2016-06-08,14,P,0.1,0.25,0.1,0.25,623,2845,0.2,1,0.35,14.24,14.24,0.2455,7302,",
		"^VIX,2016-06-01,2016-06-08,14,P,0.1,0.25,0.1,0.25,623,2845,0.2,1,0.35,14.24,14.24,0.2455,7302,",
		"^VIX,2016-06-01,2016-06-08,14.5,C,1,1.1,0.65,0.65,55,265,0.55,4233,0.7,14.24,14.24,0.7764,303,",
		"^VIX,2016-06-01,2016-06-08,14.5,P,0.3,0.55,0.3,0.55,67,1939,0.45,1893,0.65,14.24,14.24,0.4485,152,",
		"^VIX,2016-06-01,2016-06-08,15,C,0.9,0.95,0.5,0.5,82,3698,0.35,221,0.5,14.24,14.24,0.6689,1208,",
		"",
	}
	expected := strings.Join(expectedRows, "\n")
	data := string(b.Bytes())
	if data != expected {
		t.Errorf("s.Display() = \n%+v, expected \n%+v len, %d vs %d",
			[]rune(data),
			[]rune(expected),
			len(data),
			len(expected))
	}
}
