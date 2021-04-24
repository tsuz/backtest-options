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
		"underlying_symbol,quote_date,expiration,strike,option_type,open,high,low,close,trade_volume,bid_size_eod,bid_eod,ask_size_eod,ask_eod,underlying_bid_eod,underlying_ask_eod,vwap,open_interest,delivery_code",
		"^VIX,2016-06-03,2016-06-08,14,P,0.1,0.25,0.1,0.25,623,10,0.25,271,0.3,14.2,14.2,0.2455,7302,",
		"^VIX,2016-06-01,2016-06-08,14,P,0.1,0.25,0.1,0.25,623,10,0.25,271,0.3,14.2,14.2,0.2455,7302,",
		"^VIX,2016-06-01,2016-06-08,14.5,C,1,1.1,0.65,0.65,55,3499,0.55,4249,0.75,14.2,14.2,0.7764,303,",
		"^VIX,2016-06-01,2016-06-08,14.5,P,0.3,0.55,0.3,0.55,67,4331,0.4,2690,0.6,14.2,14.2,0.4485,152,",
		"^VIX,2016-06-01,2016-06-08,15,C,0.9,0.95,0.5,0.5,82,6722,0.35,3932,0.55,14.2,14.2,0.6689,1208,",
		"",
	}
	expected := strings.Join(expectedRows, "\n")
	data := string(b.Bytes())
	if data != expected {
		t.Errorf("s.Display() = \n%+v, expected \n%+v len, %d vs %d",
			data,
			expected,
			len(data),
			len(expected))
	}
}
