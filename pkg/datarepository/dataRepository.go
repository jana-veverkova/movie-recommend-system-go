package datarepository

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/csvhandler"
	"github.com/pkg/errors"
)

var castSourceUrl = "data/processed/cast.csv"

type DataSet struct {
	Ratings []*Rating
	Movies  map[string]*Movie
}

type Movie struct {
	MovieId  string
	Year     int
	Title    string
	Genres   []string
	Director string
	Actor1   string
	Actor2   string
}

type Rating struct {
	MovieId string
	UserId  string
	Value   float32
}

type Cast struct {
	MovieId  string
	Director string
	Actor1   string
	Actor2   string
}

func GetData(dataSourceUrl string) *DataSet {
	var wgSender sync.WaitGroup
	var wgReceiver sync.WaitGroup
	var wgStructCreator sync.WaitGroup

	chRecordsMovies := make(chan []string)
	chRecordsRating := make(chan []string)
	chMovies := make(chan *Movie)
	chRatings := make(chan *Rating)
	chCasts := make(chan *Cast)
	chQuit := make(chan bool)

	chRecords := csvhandler.ReadCsvData(dataSourceUrl, false, &wgSender)
	chRecordsCast := csvhandler.ReadCsvData(castSourceUrl, true, &wgSender)

	chOutput := createDataSet(chMovies, chRatings, chCasts, chQuit)
	wgReceiver.Add(1)
	go receiveRecord(chRecords, chRecordsMovies, chRecordsRating, &wgReceiver)

	wgStructCreator.Add(3)
	go processMovie(chRecordsMovies, chMovies, &wgStructCreator)
	go processRating(chRecordsRating, chRatings, &wgStructCreator)
	go processCast(chRecordsCast, chCasts, &wgStructCreator)

	wgSender.Wait()
	wgReceiver.Wait()
	close(chRecordsMovies)
	close(chRecordsRating)
	wgStructCreator.Wait()
	close(chMovies)
	close(chCasts)
	close(chRatings)

	chQuit <- true
	dataSet := <-chOutput

	return dataSet
}

func receiveRecord(chRecords <-chan []string, chRecordsMovies chan []string, chRecordsRating chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	// receives record from csvReader and sends record to movies and rating processors
	for record := range chRecords {
		chRecordsMovies <- record
		chRecordsRating <- record
	}
}

func createDataSet(chMovies chan *Movie, chRatings chan *Rating, chCasts chan *Cast, chQuit chan bool) chan *DataSet {
	dataSet := DataSet{Movies: make(map[string]*Movie), Ratings: make([]*Rating, 0)}

	chOutput := make(chan *DataSet)

	go func() {
		for {
			select {
			case movie := <-chMovies:
				if movie == nil {
					continue
				}
				if _, ok := dataSet.Movies[movie.MovieId]; ok {
					dataSet.Movies[movie.MovieId].Year = movie.Year
					dataSet.Movies[movie.MovieId].Title = movie.Title
					dataSet.Movies[movie.MovieId].Genres = movie.Genres
				} else {
					dataSet.Movies[movie.MovieId] = &Movie{MovieId: movie.MovieId, Year: movie.Year, Title: movie.Title, Genres: movie.Genres}
				}
			case rating := <-chRatings:
				if rating == nil {
					continue
				}
				dataSet.Ratings = append(dataSet.Ratings, rating)
			case cast := <-chCasts:
				if cast == nil {
					continue
				}
				if _, ok := dataSet.Movies[cast.MovieId]; ok {
					dataSet.Movies[cast.MovieId].Director = cast.Director
					dataSet.Movies[cast.MovieId].Actor1 = cast.Actor1
					dataSet.Movies[cast.MovieId].Actor2 = cast.Actor2
				} else {
					dataSet.Movies[cast.MovieId] = &Movie{MovieId: cast.MovieId, Director: cast.Director, Actor1: cast.Actor1, Actor2: cast.Actor2}
				}
			case <-chQuit:
				chOutput <- &dataSet
			}
		}
	}()

	return chOutput
}

func processMovie(chRecordsMovies chan []string, chMovies chan *Movie, wg *sync.WaitGroup) {
	defer wg.Done()

	for record := range chRecordsMovies {
		movie, err := getMovie(record)
		if err != nil {
			fmt.Printf("Couldn't get movie from \"%s\"", record)
			fmt.Println(errors.WithStack(err))
		} else if movie != nil {
			chMovies <- movie
		}
	}
}

func processRating(chRecordsRating chan []string, chRatings chan *Rating, wg *sync.WaitGroup) {
	defer wg.Done()

	for record := range chRecordsRating {
		rating, err := getRating(record)
		if err != nil {
			fmt.Printf("Couldn't get rating from \"%s\"", record)
			fmt.Println(errors.WithStack(err))
		} else if rating != nil {
			chRatings <- rating
		}
	}
}

func processCast(chRecordsCast <-chan []string, chCasts chan *Cast, wg *sync.WaitGroup) {
	defer wg.Done()

	for record := range chRecordsCast {
		cast := getCast(record)
		chCasts <- cast
	}
}

func getMovie(row []string) (*Movie, error) {
	//movieId := row[1]
	//title := row[4][:len(row[4])-7]
	//yearStr := row[4][len(row[4])-6:][1:5]
	//genres := row[5]

	y, err := strconv.Atoi(row[4][len(row[4])-6:][1:5])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	movie := Movie{
		MovieId: row[1],
		Title:   row[4][:len(row[4])-7],
		Year:    int(y),
		Genres:  strings.Split(row[5], "|"),
	}

	return &movie, nil
}

func getRating(row []string) (*Rating, error) {
	//userId := row[0]
	//movieId := row[1]
	//ratingStr := row[2]

	s, err := strconv.ParseFloat(row[2], 32)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rating := Rating{
		MovieId: row[1],
		UserId:  row[0],
		Value:   float32(s),
	}

	return &rating, nil
}

func getCast(row []string) *Cast {
	cast := Cast{
		MovieId:  row[0],
		Director: row[1],
		Actor1:   row[2],
		Actor2:   row[3],
	}

	return &cast
}
