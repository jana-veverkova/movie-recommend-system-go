package modelevaluation

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

type Summary struct {
	Rmse float32
}

func Summarize(testSet map[string]datarepository.Rating, predictions map[string]float32) (*Summary, error) {
	// returns all statistics

	rmse, err := calculateRmse(testSet, predictions)
	if err != nil {
		return nil,  errors.WithStack(err)
	}

	summary := Summary{
		Rmse: rmse,
	}

	return &summary, nil
}