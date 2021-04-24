package strategy

import (
	"backtest-options/model"
	"bytes"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

func TestCoveredCall(t *testing.T) {
	june1, _ := time.Parse(model.DateLayout, "2006-06-01")
	july2, _ := time.Parse(model.DateLayout, "2006-07-02")
	aug2, _ := time.Parse(model.DateLayout, "2006-08-02")

	v1, _ := model.NewOHLCV(june1, "SPY", july2, "116", model.Call, "1.1", "1.1", "1.1", "1.1", "623", "1.1", "0.9", "115.5", "116.5")
	v2, _ := model.NewOHLCV(july2, "SPY", july2, "116", model.Call, "0", "0", "0", "0", "623", "1", "1", "117.5", "118.5")
	v3, _ := model.NewOHLCV(july2, "SPY", aug2, "118", model.Call, "1.1", "1.1", "1.1", "1.1", "55", "1.2", "1.1", "117.5", "118.5")
	v4, _ := model.NewOHLCV(aug2, "SPY", aug2, "118", model.Call, "0.0", "0.0", "0.0", "0.0", "55", "1", "1", "119.8", "120")

	opts := model.StrategyOpts{
		StartDate:  june1,
		MinExpDays: 28,
	}
	exp := &model.StrategyResult{
		Opts: opts,
		Execs: []model.ExecLegs{
			model.ExecLegs{
				TotalProfit: decimal.NewFromFloat(100.0),
				Leg: map[string]*model.ExecOpenClose{
					"buy-stock": &model.ExecOpenClose{
						Product: model.Stock,
						Open: model.Exec{
							Date: june1,
							Px:   decimal.NewFromInt(116),
							Qty:  decimal.NewFromInt(100),
							Side: model.Buy,
						},
						Close: model.Exec{
							Date: july2,
							Px:   decimal.NewFromInt(116), // the upside is capped at the strike value
							Qty:  decimal.NewFromInt(100),
							Side: model.Sell,
						},
					},
					"covered-call": &model.ExecOpenClose{
						Product: model.Option,
						Open: model.Exec{
							Date: june1,
							Px:   decimal.NewFromInt(1),
							Qty:  decimal.NewFromInt(1),
							Side: model.Sell,
						},
						Close: model.Exec{
							Date: july2,
							Px:   decimal.NewFromFloat(0),
							Qty:  decimal.NewFromInt(1),
							Side: model.Buy,
						},
					},
				},
			},
			model.ExecLegs{
				TotalProfit: decimal.NewFromFloat(115.0),
				Leg: map[string]*model.ExecOpenClose{
					"buy-stock": &model.ExecOpenClose{
						Product: model.Stock,
						Open: model.Exec{
							Date: july2,
							Px:   decimal.NewFromInt(118),
							Qty:  decimal.NewFromInt(100),
							Side: model.Buy,
						},
						Close: model.Exec{
							Date: aug2,
							Px:   decimal.NewFromInt(118), // the upside is capped at the strike value
							Qty:  decimal.NewFromInt(100),
							Side: model.Sell,
						},
					},
					"covered-call": &model.ExecOpenClose{
						Product: model.Option,
						Open: model.Exec{
							Date: july2,
							Px:   decimal.NewFromFloat(1.15),
							Qty:  decimal.NewFromInt(1),
							Side: model.Sell,
						},
						Close: model.Exec{
							Date: aug2,
							Px:   decimal.NewFromFloat(0.0),
							Qty:  decimal.NewFromInt(1),
							Side: model.Buy,
						},
					},
				},
			},
		},
		Meta: model.StrategyMeta{
			TotalProfit:     decimal.NewFromFloat(215.0),
			TotalExecutions: 2,
		},
	}

	testData := []model.OHLCV{v1, v2, v3, v4}

	chain, err := model.NewOptionChain(testData)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating NewOptionChain"))
	}
	st, err := NewCoveredCallStrategy(chain)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating new strategy"))
	}
	strat, err := st.Run(opts)
	if err != nil {
		t.Error(errors.Wrap(err, "Error from calling covered call"))
	}
	if len(strat.Execs) != 2 {
		t.Error(errors.Errorf("Expected %+v executions but got %+v", 2, len(strat.Execs)))
	}

	if exp.Opts.ExecMethod != strat.Opts.ExecMethod {
		t.Errorf("Expected ExecMethod to be %+v but got %+v", exp.Opts.ExecMethod, strat.Opts.ExecMethod)
	}
	if exp.Opts.MinExpDays != strat.Opts.MinExpDays {
		t.Errorf("Expected MinExpDays to be %+v but got %+v", exp.Opts.MinExpDays, strat.Opts.MinExpDays)
	}
	if exp.Opts.EndDate != strat.Opts.EndDate {
		t.Errorf("Expected EndDate to be %+v but got %+v", exp.Opts.EndDate, strat.Opts.EndDate)
	}
	if !exp.Opts.StartDate.Equal(strat.Opts.StartDate) {
		t.Errorf("Expected StartDate to be %+v but got %+v", exp.Opts.StartDate, strat.Opts.StartDate)
	}

	// test the execs
	for idx, e := range exp.Execs {
		if e.TotalProfit.String() != strat.Execs[idx].TotalProfit.String() {
			t.Errorf("Expected TotalProfit to be %+v but got %+v", e.TotalProfit, strat.Execs[idx].TotalProfit)
		}
		for k, v := range e.Leg {
			out := strat.Execs[idx].Leg[k]
			if v.Product != out.Product {
				t.Errorf("Expected Product to be %+v but got %+v",
					v.Product,
					strat.Execs[idx].Leg[k].Product)
			}
			if !out.Open.Px.Equal(v.Open.Px) {
				t.Errorf("Expected open px %+v and %+v to be the same", v.Open.Px, out.Open.Px)
			}
			if !out.Open.Qty.Equal(v.Open.Qty) {
				t.Errorf("Expected open qty %+v and %+v to be the same", v.Open.Qty, out.Open.Qty)
			}
			if out.Open.Side != v.Open.Side {
				t.Errorf("Expected open side %+v and %+v to be the same", v.Open.Side, out.Open.Side)
			}
			if !out.Close.Date.Equal(v.Close.Date) {
				t.Errorf("Expected close date %+v and %+v to be the same", v.Close.Date, out.Close.Date)
			}
			if !out.Close.Px.Equal(v.Close.Px) {
				t.Errorf("Expected close px %+v and %+v to be the same for %+v", v.Close.Px, out.Close.Px, out)
			}
			if out.Close.Side != v.Close.Side {
				t.Errorf("Expected close side %+v and %+v to be the same for %+v", v.Close.Side, out.Close.Side, out)
			}
		}
	}

	// test the meta data
	if exp.Meta.TotalProfit.String() != strat.Meta.TotalProfit.String() {
		t.Errorf("Expected total profit to be %+v but got %+v",
			exp.Meta.TotalProfit.String(),
			strat.Meta.TotalProfit.String())
	}
	if exp.Meta.TotalExecutions != strat.Meta.TotalExecutions {
		t.Errorf("Expected TotalExecutions to be %+v but got %+v",
			exp.Meta.TotalExecutions,
			strat.Meta.TotalExecutions)
	}

	var metaBuf bytes.Buffer
	if err := st.OutputMeta(&metaBuf, strat); err != nil {
		t.Error(errors.Wrap(err, "expected no error to occur when GenerateResults is ran"))
	}

	metawant := `+--------------+------------------+
| TOTAL PROFIT | TOTAL EXECUTIONS |
+--------------+------------------+
|          215 |                2 |
+--------------+------------------+
`
	if metaBuf.String() != metawant {
		t.Errorf("Expected to write %+v but got %+v",
			metawant,
			metaBuf.String())
	}

	var detailBuf bytes.Buffer
	if err := st.OutputDetail(&detailBuf, strat); err != nil {
		t.Error(errors.Wrap(err, "expected no error to occur when GenerateResults is ran"))
	}

	want := `+------------+------------+------------------+--------------+----------------+-----------------+---------------+----------------+
| OPEN DATE  | CLOSE DATE |   CALL PRODUCT   | TOTAL PROFIT | OPTION OPEN PX | OPTION CLOSE PX | STOCK OPEN PX | STOCK CLOSE PX |
+------------+------------+------------------+--------------+----------------+-----------------+---------------+----------------+
| 2006-06-01 | 2006-07-02 | 116 C 2006-07-02 |          100 |              1 |               0 |           116 |            116 |
| 2006-07-02 | 2006-08-02 | 118 C 2006-08-02 |          115 |           1.15 |               0 |           118 |            118 |
+------------+------------+------------------+--------------+----------------+-----------------+---------------+----------------+
`

	if detailBuf.String() != want {
		t.Errorf("Expected to write %+v but got %+v",
			want,
			detailBuf.String())
	}

}
