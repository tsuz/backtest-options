package cmd

import (
	"backtest-options/util"
	"encoding/csv"
	"time"

	"fmt"
	"os"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "app"}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import imports external data source into a normalized file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires import dir argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Running import on livevol data from dir: %+v", args[0])
		importDir := args[0]
		outputDir := "./data"
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			log.Fatal(errors.Wrapf(err, "Error making dir at %+v", importDir))
		}
		now := time.Now().Format(time.RFC3339)
		file, err := os.OpenFile(outputDir+"/"+now+".csv", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Error opening options.csv"))
		}
		w := csv.NewWriter(file)
		if err := util.NewLiveVolImporter(w).ImportFolder(importDir); err != nil {
			log.Error(errors.Wrapf(err, "Error importing from folder %+v", importDir))
		}
		w.Flush()
		log.Info("Successfully finished importing")
	},
}

// Execute is the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
}
