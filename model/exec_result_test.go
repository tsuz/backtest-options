package model

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func TestNewOpenExec(t *testing.T) {
	now := time.Now()
	px := decimal.NewFromFloat(165.23)
	qty := decimal.NewFromInt(100)
	side := Buy
	name := "test1"
	exec := NewOpenExec(Stock, now, px, qty, side, name)
	if !exec.Open.Date.Equal(now) {
		t.Error(errors.Errorf("Expected date to be %+v but got %+v", exec.Open.Date, now))
	}
	if !exec.Open.Px.Equal(px) {
		t.Error(errors.Errorf("Expected price to be %+v but got %+v", exec.Open.Px, px))
	}
	if exec.Open.Side != side {
		t.Error(errors.Errorf("Expected side to be %+v but got %+v", exec.Open.Side, side))
	}

	closedate := now.Add(time.Hour)
	closepx := decimal.NewFromFloat(169.73)
	exec.CloseExec(closedate, closepx)
	if !exec.Close.Date.Equal(closedate) {
		t.Error(errors.Errorf("Expected date to be %+v but got %+v", exec.Close.Date, closedate))
	}
	if !exec.Close.Px.Equal(closepx) {
		t.Error(errors.Errorf("Expected price to be %+v but got %+v", exec.Close.Px, closepx))
	}
	if exec.Close.Side != Sell {
		t.Error(errors.Errorf("Expected side to be %+v but got %+v", exec.Close.Side, Sell))
	}
	if exec.Name != name {
		t.Error(errors.Errorf("Expected name to be %+v but got %+v", exec.Name, name))
	}
}
