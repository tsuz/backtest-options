package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Side represents either open or sell
type Side int

const (
	// Buy represents a purchase
	Buy Side = iota
	// Sell represents shorting or close to reduce if this is close side
	Sell
)

// ExecOpenClose is open and close exec
type ExecOpenClose struct {
	Open  Exec
	Close Exec
}

// Exec is a result of execution
type Exec struct {
	Date time.Time
	Px   decimal.Decimal
	Side Side
}

// NewOpenExec creates a new open execution
func NewOpenExec(date time.Time, px decimal.Decimal, side Side) ExecOpenClose {
	return ExecOpenClose{
		Open: Exec{
			Date: date,
			Px:   px,
			Side: side,
		},
	}
}

// CloseExec creates a closing execution and returns exec open and close
func (e ExecOpenClose) CloseExec(date time.Time, px decimal.Decimal) ExecOpenClose {
	side := Buy
	if e.Open.Side == Buy {
		side = Sell
	}
	return ExecOpenClose{
		Open: e.Open,
		Close: Exec{
			Date: date,
			Px:   px,
			Side: side,
		},
	}
}
