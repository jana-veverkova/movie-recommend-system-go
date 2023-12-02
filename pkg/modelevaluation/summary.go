package modelevaluation

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

type Summary struct {
	Rmse float32
	Mae  float32
}

func Summarize(testSet []*datarepository.Rating, predictions []float32) (*Summary, error) {
	// returns all statistics

	rmse, err := calculateRmse(testSet, predictions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mae, err := calculateMae(testSet, predictions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	summary := Summary{
		Rmse: rmse,
		Mae:  mae,
	}

	return &summary, nil
}
