package strategy

import (
	"backtest-options/model"
	"bytes"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

func TestPipStrategy(t *testing.T) {
	june1, _ := time.Parse(model.DateLayout, "2006-06-01")
	june8, _ := time.Parse(model.DateLayout, "2006-06-08")
	june15, _ := time.Parse(model.DateLayout, "2006-06-15")
	dec15, _ := time.Parse(model.DateLayout, "2006-12-15")

	v1, _ := model.NewOHLCV(june1, "SPY", june8, "116", model.Call, "0.0", "0.0", "0.0", "0.0", "623", "1.1", "0.9", "115.5", "116.5")
	v2, _ := model.NewOHLCV(june1, "SPY", dec15, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")
	v3, _ := model.NewOHLCV(june8, "SPY", june15, "118", model.Call, "0.0", "0.0", "0.0", "0.0", "0", "0.8", "0.6", "117.5", "118.5")
	v4, _ := model.NewOHLCV(june8, "SPY", dec15, "118", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "3.8", "3.6", "117.5", "118.5")
	v5, _ := model.NewOHLCV(june15, "SPY", dec15, "118", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "4.8", "4.6", "115.5", "116")

	opts := model.StrategyOpts{
		StartDate:  june1,
		MinExpDays: 28,
		PipOpts: &model.PipOpts{
			MinCallExpDTE: 4,
			MinPutExpDTE:  150,
		},
	}
	exp := &model.StrategyResult{
		Opts: opts,
		Execs: []model.ExecLegs{
			model.ExecLegs{
				TotalProfit: decimal.NewFromFloat(75.0),
				Leg: map[string]*model.ExecOpenClose{
					pipbuyStockLeg: &model.ExecOpenClose{
						Product: model.Stock,
						Open: model.Exec{
							Date: june1,
							Px:   decimal.NewFromInt(116),
							Qty:  decimal.NewFromInt(100),
							Side: model.Buy,
						},
						Close: model.Exec{
							Date: june8,
							Px:   decimal.NewFromInt(116), // the upside is capped at the strike value
							Qty:  decimal.NewFromInt(100),
							Side: model.Sell,
						},
					},
					pipcoveredCallLeg: &model.ExecOpenClose{
						Product: model.Option,
						Open: model.Exec{
							Date: june1,
							Px:   decimal.NewFromInt(1),
							Qty:  decimal.NewFromInt(1),
							Side: model.Sell,
						},
						Close: model.Exec{
							Date: june8,
							Px:   decimal.NewFromFloat(0),
							Qty:  decimal.NewFromInt(1),
							Side: model.Buy,
						},
					},
					pipfarput: &model.ExecOpenClose{
						Product: model.Option,
						Open: model.Exec{
							Date: june1,
							Px:   decimal.NewFromFloat(3.95),
							Qty:  decimal.NewFromInt(1),
							Side: model.Buy,
						},
						Close: model.Exec{
							Date: june8,
							Px:   decimal.NewFromFloat(3.7),
							Qty:  decimal.NewFromInt(1),
							Side: model.Sell,
						},
					},
				},
			},
			model.ExecLegs{
				TotalProfit: decimal.NewFromFloat(-55.0),
				Leg: map[string]*model.ExecOpenClose{
					pipbuyStockLeg: &model.ExecOpenClose{
						Product: model.Stock,
						Open: model.Exec{
							Date: june8,
							Px:   decimal.NewFromInt(118),
							Qty:  decimal.NewFromInt(100),
							Side: model.Buy,
						},
						Close: model.Exec{
							Date: june15,
							Px:   decimal.NewFromFloat(115.75),
							Qty:  decimal.NewFromInt(100),
							Side: model.Sell,
						},
					},
					pipcoveredCallLeg: &model.ExecOpenClose{
						Product: model.Option,
						Open: model.Exec{
							Date: june8,
							Px:   decimal.NewFromFloat(0.7),
							Qty:  decimal.NewFromInt(1),
							Side: model.Sell,
						},
						Close: model.Exec{
							Date: june15,
							Px:   decimal.NewFromFloat(0.0),
							Qty:  decimal.NewFromInt(1),
							Side: model.Buy,
						},
					},
					pipfarput: &model.ExecOpenClose{
						Product: model.Option,
						Open: model.Exec{
							Date: june8,
							Px:   decimal.NewFromFloat(3.7),
							Qty:  decimal.NewFromInt(1),
							Side: model.Buy,
						},
						Close: model.Exec{
							Date: june15,
							Px:   decimal.NewFromFloat(4.7),
							Qty:  decimal.NewFromInt(1),
							Side: model.Sell,
						},
					},
				},
			},
		},
		Meta: model.StrategyMeta{
			TotalProfit:     decimal.NewFromFloat(20.0),
			TotalExecutions: 2,
		},
	}

	testData := []model.OHLCV{v1, v2, v3, v4, v5}

	chain, err := model.NewOptionChain(testData)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating NewOptionChain"))
	}
	st, err := NewPIPStrategy(chain)
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

	metawant := `+----------------+------------------+--------------+------------------+
|  TOTAL PROFIT  | TOTAL EXECUTIONS | MAX DRAWDOWN |    BUY & HOLD    |
+----------------+------------------+--------------+------------------+
| 20.00 (0.17 %) |                2 |        73.33 | -25.00 (-0.22 %) |
+----------------+------------------+--------------+------------------+
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

	want := `+------------+------------+------------------+------------------+--------------+----------------------+-------------+--------------+------------+---------------+----------------+-------------------+
| OPEN DATE  | CLOSE DATE |   CALL PRODUCT   |   PUT PRODUCT    | TOTAL PROFIT | COVERED CALL PREMIUM | PUT OPEN PX | PUT CLOSE PX | PUT PROFIT | STOCK OPEN PX | STOCK CLOSE PX | CUMULATIVE PROFIT |
+------------+------------+------------------+------------------+--------------+----------------------+-------------+--------------+------------+---------------+----------------+-------------------+
| 2006-06-01 | 2006-06-08 | 116 C 2006-06-08 | 116 P 2006-12-15 |           75 |                    1 |        3.95 |         3.70 |      -0.25 |           116 |            116 |                75 |
| 2006-06-08 | 2006-06-15 | 118 C 2006-06-15 | 118 P 2006-12-15 |          -55 |                  0.7 |        3.70 |         4.70 |       1.00 |           118 |         115.75 |                20 |
+------------+------------+------------------+------------------+--------------+----------------------+-------------+--------------+------------+---------------+----------------+-------------------+
`

	if detailBuf.String() != want {
		t.Errorf("Expected to write %+v but got %+v",
			want,
			detailBuf.String())
	}

}

func TestPipStrategyInvalidParams(t *testing.T) {
	// invalid params for pip strategy

	testTable := []model.PipOpts{
		model.PipOpts{
			// MinCallExpDTE: 4,
			MinPutExpDTE: 150,
			TgtCallPxMul: decimal.Decimal(decimal.NewFromFloat(1.02)),
			TgtPutPxMul:  decimal.Decimal(decimal.NewFromFloat(0.95)),
		},
		model.PipOpts{
			MinCallExpDTE: 4,
			// MinPutExpDTE:  150,
			TgtCallPxMul: decimal.Decimal(decimal.NewFromFloat(1.02)),
			TgtPutPxMul:  decimal.Decimal(decimal.NewFromFloat(0.95)),
		},
		model.PipOpts{
			MinCallExpDTE: 4,
			MinPutExpDTE:  150,
			// TgtCallPxMul:   decimal.Decimal(decimal.NewFromFloat(1.02)),
			TgtPutPxMul: decimal.Decimal(decimal.NewFromFloat(0.95)),
		},
		model.PipOpts{
			MinCallExpDTE: 4,
			MinPutExpDTE:  150,
			TgtCallPxMul:  decimal.Decimal(decimal.NewFromFloat(1.02)),
			// TgtPutPxMul:    decimal.Decimal(decimal.NewFromFloat(0.95)),
		},
	}

	for i, v := range testTable {
		strat, _ := NewPIPStrategy(nil)
		err := strat.Validate(model.StrategyOpts{
			PipOpts: &v,
		})
		if err == nil {
			t.Errorf("Expected error but did not get error for idx %d", i)
		}
	}
}

func TestPipStrategyOptsMinCallExpDTE(t *testing.T) {

	june1, _ := time.Parse(model.DateLayout, "2006-06-01")
	june8, _ := time.Parse(model.DateLayout, "2006-06-08")
	june15, _ := time.Parse(model.DateLayout, "2006-06-15")
	dec15, _ := time.Parse(model.DateLayout, "2006-12-15")

	v1, _ := model.NewOHLCV(june1, "SPY", june8, "116", model.Call, "0.0", "0.0", "0.0", "0.0", "623", "1.1", "0.9", "115.5", "116.5")
	v2, _ := model.NewOHLCV(june1, "SPY", june15, "116", model.Call, "0.0", "0.0", "0.0", "0.0", "623", "1.5", "1.3", "115.5", "116.5")
	v3, _ := model.NewOHLCV(june1, "SPY", dec15, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")
	v4, _ := model.NewOHLCV(june8, "SPY", june8, "116", model.Call, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")
	v5, _ := model.NewOHLCV(june8, "SPY", dec15, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")
	v6, _ := model.NewOHLCV(june15, "SPY", june15, "116", model.Call, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")
	v7, _ := model.NewOHLCV(june15, "SPY", dec15, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")

	tt := []struct {
		calldte int
		r       model.ExecOpenClose
	}{
		{
			calldte: 7,
			r: model.ExecOpenClose{
				Product: model.Option,
				Open: model.Exec{
					Date: june1,
					Px:   decimal.NewFromInt(1),
					Qty:  decimal.NewFromInt(1),
					Side: model.Sell,
				},
				Close: model.Exec{
					Date: june8,
					Px:   decimal.NewFromFloat(0),
					Qty:  decimal.NewFromInt(1),
					Side: model.Buy,
				},
			},
		},
		{
			calldte: 14,
			r: model.ExecOpenClose{
				Product: model.Option,
				Open: model.Exec{
					Date: june1,
					Px:   decimal.NewFromFloat(1.4),
					Qty:  decimal.NewFromInt(1),
					Side: model.Sell,
				},
				Close: model.Exec{
					Date: june15,
					Px:   decimal.NewFromFloat(0),
					Qty:  decimal.NewFromInt(1),
					Side: model.Buy,
				},
			},
		},
	}

	testData := []model.OHLCV{v1, v2, v3, v4, v5, v6, v7}

	chain, err := model.NewOptionChain(testData)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating NewOptionChain"))
	}
	st, err := NewPIPStrategy(chain)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating new strategy"))
	}
	for idx, tab := range tt {
		opts := model.StrategyOpts{
			PipOpts: &model.PipOpts{
				MinCallExpDTE: tab.calldte,
				MinPutExpDTE:  150,
			},
		}
		strat, err := st.Run(opts)
		if err != nil {
			t.Error(errors.Wrap(err, "Error from calling pip"))
		}
		if len(strat.Execs) != 1 {
			t.Fatalf("Expected %+v executions but got %d, idx: %d", 1, len(strat.Execs), idx)
		}
		if !strat.Execs[0].Leg[pipcoveredCallLeg].Open.Px.Equal(tab.r.Open.Px) {
			t.Errorf("Expected price to be %+v  but got %+v at idx: %d", strat.Execs[0].Leg[pipcoveredCallLeg].Open.Px, tab.r.Open.Px, idx)
		}
		if !strat.Execs[0].Leg[pipcoveredCallLeg].Close.Date.Equal(tab.r.Close.Date) {
			t.Errorf("Expected price to be %+v  but got %+v at idx: %d", strat.Execs[0].Leg[pipcoveredCallLeg].Close.Date, tab.r.Close.Date, idx)
		}
	}
}

func TestPipStrategyOptsMinPutExpDTE(t *testing.T) {

	june1, _ := time.Parse(model.DateLayout, "2006-06-01")
	june8, _ := time.Parse(model.DateLayout, "2006-06-08")
	dec15, _ := time.Parse(model.DateLayout, "2006-12-15")
	mar1, _ := time.Parse(model.DateLayout, "2007-03-15")

	v1, _ := model.NewOHLCV(june1, "SPY", june8, "116", model.Call, "0.0", "0.0", "0.0", "0.0", "623", "1.1", "0.9", "115.5", "116.5")
	v2, _ := model.NewOHLCV(june1, "SPY", dec15, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")
	v3, _ := model.NewOHLCV(june1, "SPY", mar1, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "6.3", "6.1", "115.5", "116.5")
	v4, _ := model.NewOHLCV(june8, "SPY", june8, "116", model.Call, "0.0", "0.0", "0.0", "0.0", "0", "4.1", "3.8", "115.5", "116.5")
	v5, _ := model.NewOHLCV(june8, "SPY", dec15, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "4.0", "3.8", "115.5", "116.5")
	v6, _ := model.NewOHLCV(june8, "SPY", mar1, "116", model.Put, "0.0", "0.0", "0.0", "0.0", "0", "6.2", "6.0", "115.5", "116.5")

	tt := []struct {
		putdte int
		r      model.ExecOpenClose
	}{
		{
			putdte: 150,
			r: model.ExecOpenClose{
				Product: model.Option,
				Open: model.Exec{
					Px: decimal.NewFromFloat(3.95),
				},
				Close: model.Exec{
					Px: decimal.NewFromFloat(3.9),
				},
			},
		},
		{
			putdte: 200,
			r: model.ExecOpenClose{
				Product: model.Option,
				Open: model.Exec{
					Px: decimal.NewFromFloat(6.2),
				},
				Close: model.Exec{
					Px: decimal.NewFromFloat(6.1),
				},
			},
		},
	}

	testData := []model.OHLCV{v1, v2, v3, v4, v5, v6}

	chain, err := model.NewOptionChain(testData)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating NewOptionChain"))
	}
	st, err := NewPIPStrategy(chain)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating new strategy"))
	}
	for idx, tab := range tt {
		opts := model.StrategyOpts{
			PipOpts: &model.PipOpts{
				MinPutExpDTE:  tab.putdte,
				MinCallExpDTE: 4,
			},
		}
		strat, err := st.Run(opts)
		if err != nil {
			t.Error(errors.Wrap(err, "Error from calling pip"))
		}
		if len(strat.Execs) != 1 {
			t.Fatalf("Expected %+v executions but got %d, idx: %d", 1, len(strat.Execs), idx)
		}
		if !strat.Execs[0].Leg[pipfarput].Open.Px.Equal(tab.r.Open.Px) {
			t.Errorf("Expected price to be %+v  but got %+v at idx: %d", strat.Execs[0].Leg[pipfarput].Open.Px, tab.r.Open.Px, idx)
		}
		if !strat.Execs[0].Leg[pipfarput].Close.Px.Equal(tab.r.Close.Px) {
			t.Errorf("Expected price to be %+v  but got %+v at idx: %d", strat.Execs[0].Leg[pipfarput].Close.Px, tab.r.Close.Px, idx)
		}
	}
}
