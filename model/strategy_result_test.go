package model

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func TestStrategyResult(t *testing.T) {
	june1, _ := time.Parse(DateLayout, "2006-06-01")
	july2, _ := time.Parse(DateLayout, "2006-07-02")
	opts := StrategyOpts{
		ExecMethod: ExecMethodCrossSpread,
	}
	result := NewStrategyResult(opts)
	leg := map[string]*ExecOpenClose{
		"foo": &ExecOpenClose{
			Product: Stock,
			Open: Exec{
				Date: june1,
				Px:   decimal.NewFromInt(116),
				Qty:  decimal.NewFromInt(100),
				Side: Buy,
			},
			Close: Exec{
				Date: july2,
				Px:   decimal.NewFromFloat(118.4), // the upside is capped at the strike value
				Qty:  decimal.NewFromInt(100),
				Side: Sell,
			},
		},
	}
	exec, err := NewExecLegs(leg)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error creating exec legs"))
	}
	result.AddExec(exec)

	if result.Meta.TotalExecutions != 1 {
		t.Errorf("Expected %d executions but got %d", 1, result.Meta.TotalExecutions)
	}
	if result.Meta.TotalProfit.String() != "240" {
		t.Errorf("Expected total profit of %+v but got %+v",
			"240",
			result.Meta.TotalProfit.String())
	}

}
