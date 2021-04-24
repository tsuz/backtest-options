package strategy

import (
	"backtest-options/model"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type coveredCall struct {
	optchain *model.OptChainList
}

// NewCoveredCallStrategy is a new covererd call strategy
func NewCoveredCallStrategy(chain *model.OptChainList) (Strategy, error) {
	return &coveredCall{
		optchain: chain,
	}, nil
}

var coveredCallLeg = "covered-call"
var buyStockLeg = "buy-stock"

// CoveredCall
func (s *coveredCall) Run(opts model.StrategyOpts) (*model.StrategyResult, error) {

	newstrat := model.NewStrategyResult(opts)

	start := opts.StartDate
	minexpday := opts.MinExpDays

	for {
		optchain := s.optchain.GetOptionChainForQuoteDate(start, false)
		if optchain == nil {
			log.Warnf("Exiting since quote does not exist for date %+v", start)
			break
		}
		quotedate := optchain.QuoteDate
		px := optchain.UndPx
		expdate := quotedate.AddDate(0, 0, minexpday)
		expchain := optchain.GetOptionChainForExpiryDate(expdate, false)
		if expchain == nil {
			log.Warnf("Exiting since expire does not exist for date %+v, for quote date: %+v", expdate, start)
			break
		}
		strike := expchain.GetOptionChainForStrike(px, false)
		if strike == nil {
			log.Warnf("Exiting since strike does not exist for price %+v, expire date %+v, for quote date: %+v", px, expdate, start)
			break
		}

		// purchase 100 underlying stocks
		stkqty := decimal.NewFromInt(100)
		stkleg := model.NewOpenExec(
			model.Stock,
			quotedate,
			px,
			stkqty,
			model.Buy,
			"Stock")

		// write 1 option contract
		optqty := decimal.NewFromInt(1)
		optleg := model.NewOpenExec(
			model.Option,
			quotedate,
			strike.Call.AskBidMid,
			optqty,
			model.Sell,
			fmt.Sprintf("%+v C %+v", strike.S.String(), expchain.ExpireDate.Format("2006-01-02")),
		)

		expire := expchain.ExpireDate
		log.Debugf("Quotedate: %+v, expiredate: %+v", quotedate, expdate)
		expiredquote := s.optchain.GetOptionChainForQuoteDate(expire, false)
		if expiredquote == nil {
			log.Debugf("Exiting since GetOptionChainForQuoteDate does not exist for quotedate: %+v, expiredate: %+v, start: %+v",
				quotedate,
				expire,
				start)
			break
		}
		endpx := expiredquote.UndPx

		// close 100 underlying stocks, capped at the lower of the price and strike
		adjendpx := endpx
		if endpx.GreaterThan(strike.S) {
			adjendpx = strike.S
		}
		stkleg.CloseExec(expire, adjendpx)

		optleg.CloseExec(expire, decimal.NewFromInt(0))

		legs := map[string]*model.ExecOpenClose{
			coveredCallLeg: optleg,
			buyStockLeg:    stkleg,
		}
		execlegs, err := model.NewExecLegs(legs)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating new exec legs")
		}
		if err := newstrat.AddExec(execlegs); err != nil {
			return nil, errors.Wrapf(err, "Error adding exec for legs %+v", execlegs)
		}

		// the expiry date is the new start date
		start = expire
	}
	return newstrat, nil
}

// OutputDetail generates execution results
func (s *coveredCall) OutputDetail(w io.Writer, r *model.StrategyResult) error {

	data := [][]string{}

	for _, ex := range r.Execs {

		cc, ok := ex.Leg[coveredCallLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", coveredCallLeg)
		}
		stk, ok := ex.Leg[buyStockLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", buyStockLeg)
		}

		d := []string{
			cc.Open.Date.Format(model.DateLayout),
			cc.Close.Date.Format(model.DateLayout),
			cc.Name,
			ex.TotalProfit.String(),
			cc.Open.Px.String(),
			cc.Close.Px.String(),
			stk.Open.Px.String(),
			stk.Close.Px.String(),
		}
		data = append(data, d)
	}

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Open Date",
		"Close Date",
		"Call Product",
		"Total Profit",
		"Option Open Px",
		"Option Close Px",
		"Stock Open Px",
		"Stock Close Px",
	})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	return nil
}

// OutputMeta generates meta results
func (s *coveredCall) OutputMeta(w io.Writer, r *model.StrategyResult) error {

	data := [][]string{}

	for _, ex := range r.Execs {

		cc, ok := ex.Leg[coveredCallLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", coveredCallLeg)
		}
		stk, ok := ex.Leg[buyStockLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", buyStockLeg)
		}

		d := []string{
			cc.Open.Date.Format(model.DateLayout),
			cc.Close.Date.Format(model.DateLayout),
			ex.TotalProfit.String(),
			cc.Open.Px.String(),
			cc.Close.Px.String(),
			stk.Open.Px.String(),
			stk.Close.Px.String(),
		}
		data = append(data, d)
	}

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Total Profit",
		"Total Executions",
	})
	data = [][]string{
		[]string{
			r.Meta.TotalProfit.String(),
			fmt.Sprintf("%d", r.Meta.TotalExecutions),
		},
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	return nil
}
