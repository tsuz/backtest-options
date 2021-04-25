package strategy

import (
	"backtest-options/model"
	"io"
)

// Strategy is a strategy interface
type Strategy interface {
	Run(opts model.StrategyOpts) (*model.StrategyResult, error)
	OutputMeta(w io.Writer, s *model.StrategyResult) error
	OutputDetail(w io.Writer, s *model.StrategyResult) error
	Validate(opts model.StrategyOpts) error
}
