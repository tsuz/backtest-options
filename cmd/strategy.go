package cmd

import (
	"backtest-options/model"
	"backtest-options/strategy"
	"backtest-options/util"
	"encoding/csv"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	cobra "github.com/spf13/cobra"
)

var strategyCmd = &cobra.Command{
	Use:   "strategy",
	Short: "strategy runs a strategy on data",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires strategy argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Running strategy: %+v", args[0])
		dataDir := "./data"
		opts := model.StrategyOpts{
			ExecMethod: model.ExecMethodCrossSpread,
			MinExpDays: 28,
			StartDate:  time.Time{},
		}
		files, err := ioutil.ReadDir(dataDir)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "Error reading dir: %+v", dataDir))
		}
		log.Infof("Reading files from %+v", dataDir)
		ohlcvs := make([]model.OHLCV, 0)
		for _, file := range files {
			f, err := os.Open(dataDir + "/" + file.Name())
			if err != nil {
				log.Fatal(errors.Wrap(err, "Error opening options.csv"))
			}
			csvr := csv.NewReader(f)
			ohlcv, err := util.NewFileReader().ReadNormalizedCSVFile(csvr)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Error opening options.csv"))
			}
			ohlcvs = append(ohlcvs, ohlcv...)
		}
		if len(ohlcvs) == 0 {
			log.Fatal("Could not find any valid data")
		}
		log.Infof("Generated options chain")

		chain, err := model.NewOptionChain(ohlcvs)
		if err != nil {
			log.Fatal("Failed to make option chain")
		}

		switch args[0] {
		case "coveredcall":
			cc(chain, opts)
			break
		case "pip":
			pip(chain, opts)
			break
		}

		log.Info("Successfully finished running")
	},
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
