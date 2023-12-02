package moviescast

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/csvhandler"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func DownloadData(moviesDataSourceUrl string, targetUrl string) error {
	sourceUrl := "https://www.csfd.cz/hledat/?q="

	data := datarepository.GetData(moviesDataSourceUrl)

	fmt.Println("data loaded")

	var wgReceiver sync.WaitGroup
	var wgSender sync.WaitGroup

	chWriter := csvhandler.WriteCsvData(targetUrl, false, &wgReceiver)

	counter := 0
	for _, movie := range data.Movies {
		if movie.Director != "" {
			continue
		}
		if counter == 10000 {
			break
		}

		wgSender.Add(1)
		go readHtmlTable(sourceUrl, chWriter, &wgSender, *movie)

		counter++
		fmt.Print(movie.MovieId, " ")
		time.Sleep(2 * time.Second)
	}

	wgSender.Wait()
	close(chWriter)
	wgReceiver.Wait()

	return nil
}

func readHtmlTable(sourceUrl string, ch chan []string, wg *sync.WaitGroup, movie datarepository.Movie) {
	defer wg.Done()

	movieId := movie.MovieId
	title := movie.Title

	response, err := http.Get(sourceUrl + url.QueryEscape(title+" "+fmt.Sprint(movie.Year)))

	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	z := html.NewTokenizer(strings.NewReader(string(body)))

	director := ""
	actor1 := ""
	actor2 := ""

	isCreatorP := false
	isA := false
	for {
		tt := z.Next()

		tag, hasAttr := z.TagName()

		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return
			}
			fmt.Printf("Error: %v", z.Err())
			return
		}

		if string(tag) == "" {
			if isA && isCreatorP {
				if director == "" {
					director = z.Token().Data
				} else if actor1 == "" {
					actor1 = z.Token().Data
				} else if actor2 == "" {
					actor2 = z.Token().Data
				} else {
					ch <- removeNonprintables([]string{movieId, director, actor1, actor2})
					break
				}
			}
		}

		if string(tag) == "p" {
			if hasAttr {
				for {
					attrKey, attrValue, moreAttr := z.TagAttr()
					if string(attrKey) == "class" && string(attrValue) == "film-creators" {
						if z.Token().Type == html.StartTagToken {
							isCreatorP = true
						} else {
							isCreatorP = false
						}
						break
					}
					if !moreAttr {
						break
					}
				}
			}
		} else if string(tag) == "a" && z.Token().Type == html.StartTagToken {
			isA = true
		} else if string(tag) == "a" && z.Token().Type == html.EndTagToken {
			isA = false
		}
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
