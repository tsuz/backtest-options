package strategy

import (
	"backtest-options/model"
	"fmt"
	"io"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type pip struct {
	optchain *model.OptChainList
}

// NewPIPStrategy is a new covererd PIP strategy
func NewPIPStrategy(chain *model.OptChainList) (Strategy, error) {
	return &pip{
		optchain: chain,
	}, nil
}

var pipcoveredCallLeg = "covered-call"
var pipbuyStockLeg = "buy-stock"
var pipfarput = "far-put"

// CoveredCall
func (s *pip) Run(opts model.StrategyOpts) (*model.StrategyResult, error) {

	newstrat := model.NewStrategyResult(opts)

	start := opts.StartDate
	shortCallMinDays := 4
	longPutMinDays := 160

	for {
		optchain := s.optchain.GetOptionChainForQuoteDate(start, false)
		if optchain == nil {
			log.Warnf("Exiting since quote does not exist for date %+v", start)
			break
		}
		quotedate := optchain.QuoteDate
		px := optchain.UndPx
		callpx := px
		callexpdate := quotedate.AddDate(0, 0, shortCallMinDays)
		callstrike := s.getStrikePx(optchain, callexpdate, callpx)
		if callstrike == nil {
			log.Warnf("Exiting since strike does not exist for price %+v, expire date %+v, for quote date: %+v", callpx, callexpdate, start)
			break
		}

		putexpdate := quotedate.AddDate(0, 0, longPutMinDays)
		putpx := px
		putstrike := s.getStrikePx(optchain, putexpdate, putpx)
		if putstrike == nil {
			log.Warnf("Exiting since strike does not exist for price %+v, expire date %+v, for quote date: %+v", putpx, putexpdate, start)
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
			callstrike.Call.AskBidMid,
			optqty,
			model.Sell,
			fmt.Sprintf("%+v C %+v", callstrike.S.String(), callstrike.Exp.Format("2006-01-02")),
		)

		putleg := model.NewOpenExec(
			model.Option,
			quotedate,
			putstrike.Put.AskBidMid,
			optqty,
			model.Buy,
			fmt.Sprintf("%+v C %+v", putstrike.S.String(), putstrike.Exp.Format("2006-01-02")),
		)

		expire := callstrike.Exp
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
		if endpx.GreaterThan(callstrike.S) {
			adjendpx = callstrike.S
		}

		stkleg.CloseExec(expire, adjendpx)

		optleg.CloseExec(expire, decimal.NewFromInt(0))

		putendstrike := s.getStrikePx(expiredquote, putstrike.Exp, putstrike.S)
		if putendstrike == nil {
			log.Warnf("Exiting since strike does not exist for price %+v, expire date %+v, for quote date: %+v", putstrike.S, expire, expiredquote)
			break
		}
		putleg.CloseExec(expire, putendstrike.Put.AskBidMid)

		legs := map[string]*model.ExecOpenClose{
			pipcoveredCallLeg: optleg,
			pipbuyStockLeg:    stkleg,
			pipfarput:         putleg,
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
		// if expire.Weekday() == time.Saturday {
		// 	start = expire.Add(time.Hour * 24 * -1)
		// }
	}
	return newstrat, nil
}

// getStrikePx returns strike price for a given option chain, expire time and nearest price
func (s *pip) getStrikePx(optchain *model.OptChain, expd time.Time, px decimal.Decimal) *model.OptChainStrike {
	chain := optchain.GetOptionChainForExpiryDate(expd, false)
	if chain == nil {
		log.Warnf("Exiting since expire does not exist for expire date %+v", expd)
		return nil
	}
	strike := chain.GetOptionChainForStrike(px, false)
	if strike == nil {
		log.Warnf("Exiting since strike does not exist for price %+v, expire date %+v", px, expd)
		return nil
	}
	return strike
}

// OutputDetail generates execution results
func (s *pip) OutputDetail(w io.Writer, r *model.StrategyResult) error {

	data := [][]string{}
	cumprofit := decimal.Decimal{}

	for _, ex := range r.Execs {

		cc, ok := ex.Leg[pipcoveredCallLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", pipcoveredCallLeg)
		}
		stk, ok := ex.Leg[pipbuyStockLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", pipbuyStockLeg)
		}
		put, ok := ex.Leg[pipfarput]
		if !ok {
			return errors.Errorf("Error %+v key is not included", pipfarput)
		}

		cumprofit = cumprofit.Add(ex.TotalProfit)
		d := []string{
			cc.Open.Date.Format(model.DateLayout),
			cc.Close.Date.Format(model.DateLayout),
			cc.Name,
			put.Name,
			ex.TotalProfit.String(),
			cc.Open.Px.String(),
			put.Open.Px.StringFixed(2),
			put.Close.Px.StringFixed(2),
			put.Close.Px.Sub(put.Open.Px).StringFixed(2),
			stk.Open.Px.String(),
			stk.Close.Px.String(),
			cumprofit.String(),
		}
		data = append(data, d)
	}

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Open Date",
		"Close Date",
		"Call Product",
		"Put Product",
		"Total Profit",
		"Covered Call Premium",
		"Put Open Px",
		"Put Close Px",
		"Put Profit",
		"Stock Open Px",
		"Stock Close Px",
		"Cumulative Profit",
	})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	return nil
}

// OutputMeta generates meta results
func (s *pip) OutputMeta(w io.Writer, r *model.StrategyResult) error {

	data := [][]string{}
	zero := decimal.NewFromInt(0)
	hundred := decimal.NewFromInt(100)
	cumprofit := decimal.NewFromInt(0)
	maxdrawdown := decimal.NewFromInt(0)
	one := decimal.NewFromInt(1)
	firstPx := decimal.NewFromInt(0)
	lastPx := decimal.NewFromInt(0)

	for idx, ex := range r.Execs {
		cc, ok := ex.Leg[pipcoveredCallLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", pipcoveredCallLeg)
		}
		stk, ok := ex.Leg[pipbuyStockLeg]
		if !ok {
			return errors.Errorf("Error %+v key is not included", pipbuyStockLeg)
		}
		if idx == 0 {
			firstPx = stk.Open.Px
		}
		if idx == len(r.Execs)-1 {
			lastPx = stk.Close.Px
		}

		newprofit := cumprofit.Add(ex.TotalProfit)
		if cumprofit.GreaterThan(zero) {
			drawdown := newprofit.Div(cumprofit)
			if drawdown.LessThan(one) {
				diff := one.Sub(drawdown)
				if diff.GreaterThan(maxdrawdown) {
					maxdrawdown = diff
				}
			}
		}

		cumprofit = newprofit
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
		"Max Drawdown",
		"Buy & Hold",
	})
	initbp := firstPx.Mul(hundred)
	data = [][]string{
		[]string{
			fmt.Sprintf("%s (%s %%)",
				r.Meta.TotalProfit.StringFixed(2),
				r.Meta.TotalProfit.Div(initbp).Mul(hundred).StringFixed(2)),
			fmt.Sprintf("%d", r.Meta.TotalExecutions),
			fmt.Sprintf("%s", maxdrawdown.Mul(hundred).StringFixed(2)),
			fmt.Sprintf("%s (%s %%)",
				lastPx.Sub(firstPx).Mul(hundred).StringFixed(2),
				lastPx.Sub(firstPx).Mul(hundred).Mul(hundred).Div(initbp).StringFixed(2)),
		},
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	return nil
}
