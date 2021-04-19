package strategy

import (
	"option-analysis/model"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

func TestCoveredCall(t *testing.T) {
	june1, _ := time.Parse(model.DateLayout, "2006-06-01")
	july2, _ := time.Parse(model.DateLayout, "2006-07-02")
	aug2, _ := time.Parse(model.DateLayout, "2006-08-02")

	v1, _ := model.NewOHLCV(june1, "SPY", july2, "116", model.Call, "1", "1", "1", "1", "623", "115.5", "116.5")
	v2, _ := model.NewOHLCV(july2, "SPY", july2, "116", model.Call, "0.5", "0.5", "0.5", "0.5", "623", "117.5", "118.5")
	v3, _ := model.NewOHLCV(july2, "SPY", aug2, "118", model.Call, "1.1", "1.1", "1.1", "1.1", "55", "117.5", "118.5")
	v4, _ := model.NewOHLCV(aug2, "SPY", aug2, "118", model.Call, "0.1", "0.1", "0.1", "0.1", "55", "119.8", "120")

	expTable := []model.ExecOpenClose{
		model.ExecOpenClose{
			Open: model.Exec{
				Date: june1,
				Px:   decimal.NewFromFloat(1.0),
				Side: model.Sell,
			},
			Close: model.Exec{
				Date: july2,
				Px:   decimal.NewFromFloat(0.0),
				Side: model.Buy,
			},
		},
		model.ExecOpenClose{
			Open: model.Exec{
				Date: july2,
				Px:   decimal.NewFromFloat(1.1),
				Side: model.Sell,
			},
			Close: model.Exec{
				Date: aug2,
				Px:   decimal.NewFromFloat(0.0),
				Side: model.Buy,
			},
		},
	}

	testData := []model.OHLCV{v1, v2, v3, v4}
	opts := model.StrategyOpts{
		StartDate:  june1,
		MinExpDays: 28,
	}
	st, err := NewStrategy(testData)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating new strategy"))
	}
	execs, err := st.CoveredCall(opts)
	if err != nil {
		t.Error(errors.Wrap(err, "Error from calling covered call"))
	}
	if len(execs) != 2 {
		t.Error(errors.Errorf("Expected %+v executions but got %+v", 2, len(execs)))
	}

	for idx, v := range expTable {
		out := execs[idx]
		if !out.Open.Date.Equal(v.Open.Date) {
			t.Errorf("Expected open datea %+v and %+v to be the same", v.Open.Date, out.Open.Date)
		}
		if !out.Open.Px.Equal(v.Open.Px) {
			t.Errorf("Expected open px %+v and %+v to be the same", v.Open.Px, out.Open.Px)
		}
		if out.Open.Side != v.Open.Side {
			t.Errorf("Expected open side %+v and %+v to be the same", v.Open.Side, out.Open.Side)
		}
		if !out.Close.Date.Equal(v.Close.Date) {
			t.Errorf("Expected close date %+v and %+v to be the same", v.Close.Date, out.Close.Date)
		}
		if !out.Close.Px.Equal(v.Close.Px) {
			t.Errorf("Expected close px %+v and %+v to be the same", v.Close.Px, out.Close.Px)
		}
		if out.Close.Side != v.Close.Side {
			t.Errorf("Expected close side %+v and %+v to be the same", v.Close.Side, out.Close.Side)
		}
	}
}
