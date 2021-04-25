package model

import (
	"time"

	"github.com/shopspring/decimal"
)

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
	PipOpts *PipOpts
}

// PipOpts is an option custom for pip strategy
type PipOpts struct {
	// MinCallExpDTE is the minimum number of DTE until the next expiry for the call option
	MinCallExpDTE int
	// MinPutExpDTE is the minimum number of DTE until the next expiry for the put option
	MinPutExpDTE int
	// TgtCallPxMul is to determine the target strike price by multiplying the multiplier (TgtCallPxMul) with the quote date's underlying price. If the target call strike is 2% above the current underlying price, then TgtCallPxMul is 1.02
	TgtCallPxMul decimal.Decimal
	// TgtPutPxMul is to determine the target strike price by multiplying the multiplier (TgtPutPxMul) with the quote date's underlying price. If the target put strike is 2% below the current underlying price, then TgtPutPxMul is 0.98
	TgtPutPxMul decimal.Decimal
}
