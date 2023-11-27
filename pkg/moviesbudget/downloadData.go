package moviesbudget

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func DownloadData(targetUrl string) error {
	sourceUrl := "https://www.the-numbers.com/movie/budgets/all"

	var wgReceiver sync.WaitGroup
	var wgSender sync.WaitGroup

	ch := make(chan []string, 100)

	wgReceiver.Add(1)
	go writer(targetUrl, ch, &wgReceiver)

	wgSender.Add(1)
	go readHtmlTable(sourceUrl, ch, &wgSender)

	wgSender.Wait()
	close(ch)
	wgReceiver.Wait()

	return nil
}

func readHtmlTable(sourceUrl string, ch chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(sourceUrl)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	z := html.NewTokenizer(strings.NewReader(string(body)))

	row := make([]string, 0)

	depth := 0
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		switch tt {
		case html.TextToken:
			if depth > 0 {
				re := regexp.MustCompile(`\r?\n`)
				inp := re.ReplaceAllString(z.Token().Data, " ")
				row = append(row, inp)
			}
		case html.StartTagToken, html.EndTagToken:
			if z.Token().Data == "tr" {
				if tt == html.StartTagToken {
					depth++
					row = nil
				} else {
					depth--
					ch <- row
				}
			}
		}
	}
}

func writer(fileName string, ch chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	csvFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	for row := range ch {
		err = csvwriter.Write(row)
		if err != nil {
			fmt.Println(errors.WithStack(err))
		}
		csvwriter.Flush()
	}
}
