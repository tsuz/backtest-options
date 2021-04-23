package model

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// StrategyResult is a strategy result
type StrategyResult struct {
	Opts  StrategyOpts
	Execs []ExecLegs
	Meta  StrategyMeta
}

// NewStrategyResult returns a new strategy results
func NewStrategyResult(opts StrategyOpts) *StrategyResult {
	return &StrategyResult{
		Execs: make([]ExecLegs, 0),
		Meta:  StrategyMeta{},
		Opts:  opts,
	}
}

// AddExec adds a new execution
func (r *StrategyResult) AddExec(exec ExecLegs) error {
	r.Execs = append(r.Execs, exec)
	// calculate profit
	r.Meta.TotalExecutions++

	profit, err := exec.GetProfit()
	if err != nil {
		return errors.Wrap(err, "Error getting profit")
	}
	r.Meta.TotalProfit = r.Meta.TotalProfit.Add(profit)
	return nil
}

// StrategyMeta is a meta data for the strategy
type StrategyMeta struct {
	TotalExecutions int
	TotalProfit     decimal.Decimal
}

// ExecLegs is a leg for each exec
type ExecLegs struct {
	TotalProfit decimal.Decimal
	Leg         map[string]*ExecOpenClose
}
