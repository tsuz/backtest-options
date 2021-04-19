package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// ExecLeg represents multiple executions as legs
type ExecLeg struct {
	OpenDate  time.Time
	CloseDate time.Time
	Profit    decimal.Decimal
	Legs      map[string]ExecOpenClose
}

// NewExecLeg creates a series of completed exec legs including profit and open/close date
func NewExecLeg(legs map[string]ExecOpenClose) ExecLeg {
	var open, close time.Time
	profit := decimal.Zero
	for _, v := range legs {
		if open.IsZero() {
			open = v.Open.Date
		}
		if close.IsZero() {
			close = v.Close.Date
		}
		diff := v.Close.Px.Sub(v.Open.Px)
		if v.Open.Side == Sell {
			diff = profit.Mul(decimal.NewFromFloat(-1.0))
		}
		profit = profit.Add(diff)
	}
	return ExecLeg{
		CloseDate: close,
		Legs:      legs,
		OpenDate:  open,
		Profit:    profit,
	}
}
