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
	"unicode"

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

	isTr := false
	isTd := false
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return
			}
			fmt.Printf("Error: %v", z.Err())
			return
		}

		switch tt {
		case html.TextToken:
			if isTr && isTd {
				re := regexp.MustCompile(`\r?\n`)
				inp := re.ReplaceAllString(z.Token().Data, " ")
				row = append(row, inp)
			}

		case html.StartTagToken:
			t := z.Token()
			switch t.Data {
			case "tr":
				isTr = true
				row = nil
			case "td":
				isTd = true
			}
			
		case html.EndTagToken:
			t := z.Token()
			switch t.Data {
			case "tr":
				isTr = false
				if row != nil {
					ch <- row
				}
			case "td":
				isTd = false
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
		err = csvwriter.Write(removeNonprintables(row))
		if err != nil {
			fmt.Println(errors.WithStack(err))
		}
		csvwriter.Flush()
	}
}

func removeNonprintables(items []string) []string {
	result := make([]string, 0)
	for _, item := range items {
		item = strings.Replace(item, "\u00a0", "", 100)
		result = append(result, strings.TrimFunc(item, func(r rune) bool {
			return !unicode.IsGraphic(r)
		}))
	}
	return result
}
