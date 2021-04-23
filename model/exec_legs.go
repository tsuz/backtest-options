package model

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// NewExecLegs creates new exec leg
func NewExecLegs(leg map[string]*ExecOpenClose) (ExecLegs, error) {
	l := ExecLegs{
		Leg: leg,
	}
	dec, err := l.GetProfit()
	if err != nil {
		return ExecLegs{}, errors.Wrap(err, "Error getting profit")
	}
	l.TotalProfit = dec
	return l, nil
}

// GetProfit get profit
func (l ExecLegs) GetProfit() (decimal.Decimal, error) {
	tot := decimal.Decimal{}
	for _, v := range l.Leg {
		prof, err := v.GetProfit()
		if err != nil {
			return decimal.Decimal{}, errors.Wrapf(err, "Error getting profit for %+v", v)
		}
		tot = tot.Add(prof)
	}
	return tot, nil
}
