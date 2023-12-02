package moviesbudget

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/csvhandler"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

func ProcessData(budgetDataSource string, moviesDataSource string) error {
	var wgSender sync.WaitGroup
	var wgReceiver sync.WaitGroup

	chRecords := csvhandler.ReadCsvData(budgetDataSource, true, &wgSender)
	
	moviesData := datarepository.GetData(moviesDataSource)
	
	fmt.Println("movies loaded")

	go receiveData(chRecords, moviesData, &wgReceiver)

	wgSender.Wait()
	wgReceiver.Wait()

	return nil
}

func receiveData(ch <- chan []string, moviesData *datarepository.DataSet, wg *sync.WaitGroup) {
	defer wg.Done()

	for row := range ch {
		title := row[2]
		budgetRaw := row[3]
		worldwideGrossRaw := row[5]

		fmt.Printf("Raw: Title: %s, Budget: %s, Worldwide: %s \n", title, budgetRaw, worldwideGrossRaw)

		budget, err := convertMoneyString(budgetRaw)
		if err != nil {
			fmt.Println(errors.WithStack(err))
		}

		worldwideGross, err := convertMoneyString(worldwideGrossRaw)
		if err != nil {
			fmt.Println(errors.WithStack(err))
		}

		movieId := searchMovieByTitle(title, moviesData.Movies)
		fmt.Printf("MovieId: %s, Title: %s, Budget: %d, Worldwide: %d \n", movieId, title, budget, worldwideGross)
	}
}

func convertMoneyString(moneyStr string) (int, error) {
	usCurrencyReg := `^\$?\-?([1-9]{1}[0-9]{0,2}(\,\d{3})*(\.\d{0,2})?|[1-9]{1}\d{0,}(\.\d{0,2})?|0(\.\d{0,2})?|(\.\d{1,2}))$|^\-?\$?([1-9]{1}\d{0,2}(\,\d{3})*(\.\d{0,2})?|[1-9]{1}\d{0,}(\.\d{0,2})?|0(\.\d{0,2})?|(\.\d{1,2}))$|^\(\$?([1-9]{1}\d{0,2}(\,\d{3})*(\.\d{0,2})?|[1-9]{1}\d{0,}(\.\d{0,2})?|0(\.\d{0,2})?|(\.\d{1,2}))\)`
	re := regexp.MustCompile(usCurrencyReg)

	matches := re.FindStringSubmatch(strings.Replace(moneyStr, " ", "", 100))
	if matches == nil {
		return 0, nil
	}

	parts := strings.Split(matches[1], ".")
	wholePart := strings.Replace(parts[0], ",", "", 100)
	num, err := strconv.Atoi(wholePart)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func searchMovieByTitle(title string, movies map[string]*datarepository.Movie) string {
	for key, movie := range movies {
		titles := getAllTitles(movie.Title)
		for _, val := range titles {
			if correctArticle(val) == title {
				return key
			}
		}
	}

	return ""
}

func getAllTitles(title string) []string {
	titles := make([]string, 0)

	parReg := `(^\S*)\s*\(([^)]+)\)`
	re := regexp.MustCompile(parReg)
	matches := re.FindStringSubmatch(title)

	if matches == nil {
		titles = append(titles, title)
	} else {
		titles = append(titles, matches[1])
		titles = append(titles, matches[2])
	}

	return titles
}

func correctArticle(title string) string {
	newTitle := ""

	commaIx := strings.LastIndex(title, ",")
	if commaIx == 0 || commaIx == len(title)-1 {
		newTitle = title
	} else {
		article := title[commaIx+1:]
		switch strings.ToLower(article) {
		case " the":
			newTitle = "The " + title[:commaIx]
		case " a":
			newTitle = "A " + title[:commaIx]
		default:
			newTitle = title
		}
	}

	return newTitle
}
