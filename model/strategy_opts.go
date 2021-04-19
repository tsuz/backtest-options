package model

import "time"

// ExecMethod is an execution method
type ExecMethod int

const (
	// ExecMethodCrossSpread is to represent an immediate fill that crosses the spread. If this is a sell, it will execute at the bid price. If this is a buy, it will execute at an ask price.
	ExecMethodCrossSpread = iota

	// ExecMethodMidpoint is to represent a fill at the midprice of ask and bid.
	ExecMethodMidpoint
)

// StrategyOpts is an argument for strategy
type StrategyOpts struct {
	// ExecMethod is an order execution method
	ExecMethod ExecMethod
	// MinExpDays is a minimum number of expiring days
	MinExpDays int
	// StartDate is the date in which the strategy starts executing
	StartDate time.Time
	// EndDate is the last date in which the strategy ends executing
	EndDate string
}
