package models

import (
	"fmt"
	"path"
	"strings"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/modelevaluation"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/persist"
	"github.com/pkg/errors"
)

type Model interface {
	ComputeParams(map[string]datarepository.Rating) (any, error)
	Predict(map[string]datarepository.Rating, string) (map[string]float32, error)
	GetName() string
}

func Train(m Model, dataSourceUrl string) error {
	data, err := datarepository.GetData(dataSourceUrl)
	if err != nil {
		return errors.WithStack(err)
	}

	params, err := m.ComputeParams(data.Ratings)
	if err != nil {
		return errors.WithStack(err)
	}

	_, file := path.Split(dataSourceUrl)
	fileName := strings.Split(file, ".")[0]

	err = persist.Save(fmt.Sprintf("data/modelParams/%s/%s.json", m.GetName(), fileName), params)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Evaluate(m Model, trainedOnUrl string, dataSourceUrl string) (*modelevaluation.Summary, error) {
	// get data for prediction
	data, err := datarepository.GetData(dataSourceUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// load parameters
	_, file := path.Split(trainedOnUrl)
	fileName := strings.Split(file, ".")[0]

	// create predictions
	predictions, err := m.Predict(data.Ratings, fileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// evaluate predictions
	summary, err := modelevaluation.Summarize(data.Ratings, predictions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return summary, nil
}
