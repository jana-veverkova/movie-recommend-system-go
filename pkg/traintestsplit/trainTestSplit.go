package traintestsplit

import (
	"math/rand"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/csvhandler"
	"github.com/pkg/errors"
)

func Split(sourceUrl string, targetDir string) error {
	records, err := csvhandler.ReadCsvData(sourceUrl, true)
	if err != nil {
		return errors.WithStack(err)
	}

	header := make([]string, 0)
	trainSet := make([][]string, 0)
	testSet := make([][]string, 0)

	for ix, row := range records {
		if ix == 0 {
			header = row
		} else if rand.Intn(10) == 0 {
			testSet = append(testSet, row)
		} else {
			trainSet = append(trainSet, row)
		}
	}

	err = csvhandler.WriteCsvData(trainSet, targetDir+"/train.csv", header)
	if err != nil {
		return errors.WithStack(err)
	}

	err = csvhandler.WriteCsvData(testSet, targetDir+"/test.csv", header)
	if err != nil {
		return errors.WithStack(err)
	}

	return err
}
