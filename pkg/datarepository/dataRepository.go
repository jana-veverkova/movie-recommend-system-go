package datarepository

import (
	"strconv"
	"strings"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/csvhandler"
	"github.com/pkg/errors"
)

type DataSet struct {
	Movies  map[string]Movie
	Ratings map[string]Rating
}

type Movie struct {
	MovieId string
	Year    int
	Title   string
	Genres  []string
}

type Rating struct {
	MovieId string
	UserId  string
	Value   float32
}

func GetData(dataSourceUrl string) (*DataSet, error) {
	data, err := csvhandler.ReadCsvData(dataSourceUrl, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dataSet, err := formatData(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return dataSet, err
}

func formatData(records [][]string) (*DataSet, error) {
	movies := make(map[string]Movie)
	ratings := make(map[string]Rating)

	for _, row := range records {
		userId := row[0]
		movieId := row[1]
		ratingStr := row[2]
		//timest := row[3]
		title := row[4][:len(row[4])-7]
		yearStr := row[4][len(row[4])-6:][1:5]
		genres := row[5]

		s, err := strconv.ParseFloat(ratingStr, 32)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		ratingVal := float32(s)
		yearVal := int(y)

		movies[movieId] = Movie{
			MovieId: movieId,
			Title:   title,
			Year:    yearVal,
			Genres:  strings.Split(genres, "|"),
		}

		ratings[userId+"/"+movieId] = Rating{
			MovieId: movieId,
			UserId:  userId,
			Value:   ratingVal,
		}
	}

	result := DataSet{Movies: movies, Ratings: ratings}

	return &result, nil
}
