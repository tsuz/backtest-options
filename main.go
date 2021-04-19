package main

import (
	"io/ioutil"
	"log"
	"backtest-options/model"
	"backtest-options/util"
	"time"

	"github.com/pkg/errors"
)

func main() {
	// reader := util.NewFileReader()
	_, err := readDir()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error readDir"))
	}

}

func readDir() ([]model.OHLCV, error) {
	dirName := "/Users/takutosuzuki/Documents/data/options/formatted"
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	ohlcv := make([]model.OHLCV, 0)

	for idx, f := range files {
		reader := util.NewFileReader()
		data, err := reader.ReadFile(dirName + "/" + f.Name())
		if err != nil {
			log.Fatal(errors.Wrapf(err, "Error reading file name %+v", f.Name()))
		}
		ohlcv = append(ohlcv, data...)

		if idx == 2 {
			break
		}
	}

	list, err := model.NewOptionChain(ohlcv)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error creating new options chain"))
	}
	if list == nil {

	}
	oneweek := time.Now().Add(time.Hour * 24 * -1 * 7)
	quotes := list.GetOptionChainForQuoteDate(oneweek, false)
	expiry := time.Now().Add(time.Hour * 24 * 7 * 28)
	expiries := quotes.GetOptionChainForExpiryDate(expiry, false)
	log.Printf("Quotes %+v", expiries)

	return nil, nil
}

type myCloser interface {
	Close() error
}

// closeFile is a helper function which streamlines closing
// with error checking on different file types.
func closeFile(f myCloser) {
	err := f.Close()
	check(err)
}

// check is a helper function which streamlines error checking
func check(e error) {
	if e != nil {
		panic(e)
	}
}
