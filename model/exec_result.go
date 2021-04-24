package model

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Side represents either open or sell
type Side int

type ProductType int

const (
	Stock ProductType = iota
	Option
)

const (
	// Buy represents a purchase
	Buy Side = iota
	// Sell represents shorting or close to reduce if this is close side
	Sell
)

// ExecOpenClose is open and close exec
type ExecOpenClose struct {
	Close   Exec
	Name    string
	Open    Exec
	Product ProductType
}

// Exec is a result of execution
type Exec struct {
	Date time.Time
	Px   decimal.Decimal
	Qty  decimal.Decimal
	Side Side
}

// NewOpenExec creates a new open execution
func NewOpenExec(
	product ProductType,
	date time.Time,
	px decimal.Decimal,
	qty decimal.Decimal,
	side Side,
	name string,
) *ExecOpenClose {
	return &ExecOpenClose{
		Name: name,
		Open: Exec{
			Date: date,
			Px:   px,
			Qty:  qty,
			Side: side,
		},
		Product: product,
	}
}

// CloseExec creates a closing execution and returns exec open and close
func (e *ExecOpenClose) CloseExec(date time.Time, px decimal.Decimal) {
	side := Buy
	if e.Open.Side == Buy {
		side = Sell
	}
	e.Close = Exec{
		Date: date,
		Px:   px,
		Side: side,
	}
}

// GetProfit returns profit for this execution
func (e ExecOpenClose) GetProfit() (decimal.Decimal, error) {
	diff := decimal.Decimal{}

	switch e.Open.Side {
	case Buy:
		diff = e.Close.Px.Sub(e.Open.Px)
		break
	case Sell:
		diff = e.Open.Px.Sub(e.Close.Px)
		break
	default:
		return diff, errors.Errorf("Unsupported open side %+v", e.Open.Side)
	}

	switch e.Product {
	case Stock:
		return diff.Mul(e.Open.Qty), nil
	case Option:
		return diff.Mul(decimal.NewFromInt(100)).Mul(e.Open.Qty), nil
	default:
		return diff, errors.Errorf("Unsupported product %+v", e.Product)
	}
}
