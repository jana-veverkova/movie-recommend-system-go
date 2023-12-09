package csvhandler

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var lock sync.Mutex

func ReadCsvData(url string, withHeader bool, wg *sync.WaitGroup) <-chan []string {
	chRecords := make(chan []string, 100)

	wg.Add(1)
	go func() {
		defer wg.Done()

		file, err := os.Open(url)
		if err != nil {
			fmt.Println(errors.WithStack(err))
			return
		}

		defer file.Close()

		reader := csv.NewReader(file)

		if !withHeader {
			_, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					return
				}
				fmt.Println(errors.WithStack(err))
				return
			}
		}

		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					close(chRecords)
					return
				}

				fmt.Println(errors.WithStack(err))
			}

			chRecords <- record
		}
	}()

	return chRecords
}

func WriteCsvData(fileName string, rewriteFile bool, wg *sync.WaitGroup) chan []string {
	chRecords := make(chan []string)

	wg.Add(1)
	go func() {
		defer wg.Done()

		lock.Lock()
		defer lock.Unlock()

		var csvFile *os.File

		if rewriteFile {
			f, err := os.Create(fileName)
			if err != nil {
				fmt.Println(errors.WithStack(err))
				return
			}
			csvFile = f
		} else {
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
			if err != nil {
				fmt.Println(errors.WithStack(err))
				return
			}
			csvFile = f
		}

		defer csvFile.Close()

		csvwriter := csv.NewWriter(csvFile)

		for row := range chRecords {
			err := csvwriter.Write(row)
			if err != nil {
				fmt.Println(errors.WithStack(err))
			}
			csvwriter.Flush()
		}
	}()

	return chRecords
}
