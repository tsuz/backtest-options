package model

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func TestExecLegs(t *testing.T) {

	covcall := NewOpenExec(
		Option,
		time.Now(),
		decimal.NewFromFloat(10.0),
		decimal.NewFromInt(1),
		Sell,
		"name1")
	covcall.CloseExec(time.Now(), decimal.NewFromFloat(3.3))

	stock := NewOpenExec(
		Stock,
		time.Now(),
		decimal.NewFromFloat(175.3),
		decimal.NewFromInt(100),
		Buy,
		"name2")
	stock.CloseExec(time.Now(), decimal.NewFromFloat(174.2))
	execs := map[string]*ExecOpenClose{
		"covered-call": covcall,
		"stock":        stock,
	}
	legs, err := NewExecLegs(execs)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error executing legs"))
	}
	profit, err := legs.GetProfit()
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error getting profit"))
	}
	if profit.String() != "560" {
		t.Errorf("Expected 560 as profit but got %+v", profit.String())
	}
}
