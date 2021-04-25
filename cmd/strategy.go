package cmd

import (
	"backtest-options/model"
	"backtest-options/strategy"
	"backtest-options/util"
	"encoding/csv"
	"io/ioutil"
	"os"
	"time"

	"github.com/shopspring/decimal"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	cobra "github.com/spf13/cobra"
)

func getStrategyCmd() *cobra.Command {
	strategyCmd := &cobra.Command{
		Use:   "strategy",
		Short: "runs a strategy on data",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires strategy argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("requires a strategy")
			}
		},
	}

	ccCmd := &cobra.Command{
		Use:   "coveredcall",
		Short: "runs a coveredcall strategy",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Starting Covered call strategy")

			chain, err := loadOHLCV()
			if err != nil {
				log.Fatal("Failed to make option chain")
			}

			opts := model.StrategyOpts{
				ExecMethod: model.ExecMethodCrossSpread,
				MinExpDays: 28,
				StartDate:  time.Time{},
			}

			cc(chain, opts)

			log.Info("Successfully finished running")
		},
	}

	pipCmd := &cobra.Command{
		Use:   "pip",
		Short: "runs a pip strategy",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Starting pip strategy")

			cexpd := cmd.Flag("minCallDTE")
			cexp, err := decimal.NewFromString(cexpd.Value.String())
			if err != nil {
				log.Fatal(errors.Wrapf(err, "Error parsing minCallDTE: %+v", cexpd.Value.String()))
			}

			pexpd := cmd.Flag("minPutDTE")
			pexp, err := decimal.NewFromString(pexpd.Value.String())
			if err != nil {
				log.Fatal(errors.Wrapf(err, "Error parsing minPutDTE: %+v", pexpd.Value.String()))
			}

			opts := model.StrategyOpts{
				ExecMethod: model.ExecMethodCrossSpread,
				MinExpDays: 28,
				StartDate:  time.Time{},
				PipOpts: &model.PipOpts{
					MinCallExpDTE: int(cexp.IntPart()),
					MinPutExpDTE:  int(pexp.IntPart()),
					TgtCallPxMul:  decimal.NewFromInt(1),
					TgtPutPxMul:   decimal.NewFromInt(1),
				},
			}

			chain, err := loadOHLCV()
			if err != nil {
				log.Fatal("Failed to make option chain")
			}

			pip(chain, opts)

			log.Info("Successfully finished running")
		},
	}
	pipCmd.Flags().String("minCallDTE", "4", "Minimum number of DTE for the call option (Default 4)")
	pipCmd.Flags().String("minPutDTE", "150", "Minimum number of DTE for the put option (Default: 150)")

	strategyCmd.AddCommand(pipCmd)
	strategyCmd.AddCommand(ccCmd)

	return strategyCmd
}

func loadOHLCV() (*model.OptChainList, error) {
	dataDir := "./data"
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading dir: %+v", dataDir)
	}
	log.Infof("Reading files from %+v", dataDir)
	ohlcvs := make([]model.OHLCV, 0)
	for _, file := range files {
		f, err := os.Open(dataDir + "/" + file.Name())
		if err != nil {
			return nil, errors.Wrap(err, "Error opening options.csv")
		}
		csvr := csv.NewReader(f)
		ohlcv, err := util.NewFileReader().ReadNormalizedCSVFile(csvr)
		if err != nil {
			return nil, errors.Wrap(err, "Error opening options.csv")
		}
		ohlcvs = append(ohlcvs, ohlcv...)
	}
	if len(ohlcvs) == 0 {
		return nil, errors.Wrap(err, "Could not find any valid data")
	}
	log.Infof("Generated options chain")

	chain, err := model.NewOptionChain(ohlcvs)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to make option chain")
	}
	return chain, nil
}

func cc(chain *model.OptChainList, opts model.StrategyOpts) {
	s, err := strategy.NewCoveredCallStrategy(chain)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error creating new strategy"))
	}
	log.Infof("Starting strategy with opts %+v", opts)
	result, err := s.Run(opts)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error running covered call strategy"))
	}
	stdout := os.Stdout
	s.OutputDetail(stdout, result)
	s.OutputMeta(stdout, result)
}

func pip(chain *model.OptChainList, opts model.StrategyOpts) {
	s, err := strategy.NewPIPStrategy(chain)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error creating new strategy"))
	}
	log.Infof("Starting strategy with opts %+v", opts)
	result, err := s.Run(opts)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error running covered call strategy"))
	}
	stdout := os.Stdout
	s.OutputDetail(stdout, result)
	s.OutputMeta(stdout, result)
}
