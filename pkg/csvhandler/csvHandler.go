package csvhandler

import (
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
)

func ReadCsvData(url string, withHeader bool) ([][]string, error) {
	file, err := os.Open(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !withHeader {
		// remove header
		records = records[1:]
	}

	return records, err
}

func WriteCsvData(records [][]string, fileName string, header []string) error {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return errors.WithStack(err)
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	err = csvwriter.Write(header)
	if err != nil {
		return errors.WithStack(err)
	}

	err = csvwriter.WriteAll(records)
	if err != nil {
		return errors.WithStack(err)
	}

	return err
}
