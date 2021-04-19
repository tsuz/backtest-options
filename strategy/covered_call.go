package strategy

import (
	"log"
	"option-analysis/model"

	"github.com/shopspring/decimal"
)

// CoveredCall
func (s *strategy) CoveredCall(opts model.StrategyOpts) ([]model.ExecOpenClose, error) {
	result := make([]model.ExecOpenClose, 0)

	start := opts.StartDate
	minexpday := opts.MinExpDays

	for {
		optchain := s.optchain.GetOptionChainForQuoteDate(start, false)
		if optchain == nil {
			log.Printf("Exiting since quote does not exist for date %+v", start)
			break
		}
		quotedate := optchain.QuoteDate
		px := optchain.UndPx
		expdate := start.AddDate(0, 0, minexpday)
		expchain := optchain.GetOptionChainForExpiryDate(expdate, false)
		if expchain == nil {
			log.Printf("Exiting since expire does not exist for date %+v, for quote date: %+v", expdate, start)
			break
		}
		strike := expchain.GetOptionChainForStrike(px, false)
		// execute buying 100 shares and selling the call at this strike
		if strike == nil {
			log.Printf("Exiting since strike does not exist for price %+v, expire date %+v, for quote date: %+v", px, expdate, start)
			break
		}
		exec := model.NewOpenExec(quotedate, strike.Call.Open, model.Sell)
		expire := expchain.ExpireDate

		expiredquote := s.optchain.GetOptionChainForQuoteDate(expire, false)
		endpx := expiredquote.UndPx
		adjendpx := px.Sub(endpx)

		if adjendpx.LessThan(decimal.Zero) {
			adjendpx = decimal.Zero
		}

		final := exec.CloseExec(expire, adjendpx)
		result = append(result, final)
		start = expire
	}
	return result, nil
}
