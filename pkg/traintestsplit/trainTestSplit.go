package traintestsplit

import (
	"math/rand"
	"sync"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/csvhandler"
)

func Split(sourceUrl string, targetDir string) {
	var wgSender sync.WaitGroup
	var wgReceiver sync.WaitGroup
	var wgWriter sync.WaitGroup

	chRecords := csvhandler.ReadCsvData(sourceUrl, true, &wgSender)
	chWriterTrain := csvhandler.WriteCsvData(targetDir+"/train.csv", true, &wgWriter)
	chWriterTest := csvhandler.WriteCsvData(targetDir+"/test.csv", true, &wgWriter)

	wgReceiver.Add(1)
	go func() {
		defer wgReceiver.Done()

		for row := range chRecords {
			if rand.Intn(5) == 0 {
				chWriterTest <- row
			} else {
				chWriterTrain <- row
			}
		}
	}()

	wgSender.Wait()
	wgReceiver.Wait()
	close(chWriterTest)
	close(chWriterTrain)
	wgWriter.Wait()
}
