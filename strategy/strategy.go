package strategy

import (
	"option-analysis/model"

	"github.com/pkg/errors"
)

// Strategy is a strategy interface
type Strategy interface {
	CoveredCall(opts model.StrategyOpts) ([]model.ExecOpenClose, error)
}

type strategy struct {
	// optchain is an option chain data structure
	optchain *model.OptChainList
	// opts is strategy options
	opts model.StrategyOpts
}

// NewStrategy creates a strategy
func NewStrategy(d []model.OHLCV) (Strategy, error) {
	chain, err := model.NewOptionChain(d)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating NewOptionChain")
	}
	return &strategy{
		optchain: chain,
	}, nil
}
