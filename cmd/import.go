package cmd

import (
	"backtest-options/util"
	"encoding/csv"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	cobra "github.com/spf13/cobra"
)

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
		defer w.Flush()
		if err := util.NewLiveVolImporter(w).ImportFolder(importDir); err != nil {
			log.Fatal(errors.Wrapf(err, "Error importing from folder %+v", importDir))
		}
		log.Info("Successfully finished importing")
	},
}
