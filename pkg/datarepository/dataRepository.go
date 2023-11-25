package datarepository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/csvhandler"
	"github.com/pkg/errors"
)

type dataSet struct {
	movies  map[string]movie
	ratings map[string]rating
}

type movie struct {
	movieId string
	year int
	title   string
	genres  []string
}

type rating struct {
	movieId string
	userId  string
	rating  float32
}

func GetData(dataSource string) (*dataSet, error) {
	dataSourceUrl := dataSource
	
	switch dataSource {
    case "train":
        dataSourceUrl = "data/trainTest/train.csv"
    case "test":
        dataSourceUrl = "data/trainTest/test.csv"
    case "edx":
        dataSourceUrl = "data/processed/edx.csv"
	case "holdout_test":
		dataSourceUrl = "data/processed/final_holdout_test.csv"
    }

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

func formatData(records [][]string) (*dataSet, error) {
	movies := make(map[string]movie)
	ratings := make(map[string]rating)

	for _, row := range records {
		fmt.Println(row)
		userId := row[0]
		movieId := row[1]
		ratingStr := row[2]
		//timest := row[3]
		title := row[4][:len(row[4])-7]
		yearStr := row[4][len(row[4])-6:][1:5]
		genres := row[5]

		s, err := strconv.ParseFloat(ratingStr, 32)
		if err != nil {
			fmt.Printf("String %s cannot be converted to float.", ratingStr)
			return nil, errors.WithStack(err)
		}

		y, err := strconv.Atoi(yearStr)
		if err != nil {
			fmt.Printf("String %s cannot be converted to int.", yearStr)
			return nil, errors.WithStack(err)
		}

		ratingVal := float32(s)
		yearVal := int(y)

		movies[movieId] = movie{
			movieId: movieId,
			title:   title,
			year: yearVal,
			genres:  strings.Split(genres, "|"),
		}

		ratings[userId+"/"+movieId] = rating{
			movieId: movieId,
			userId:  userId,
			rating:  ratingVal,
		}
	}

	result := dataSet{movies: movies, ratings: ratings}

	return &result, nil
}
