package models

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

type Intercept struct {
	model IModel
	intercept float32
}

func (i *Intercept) predict(data *datarepository.DataSet) ([]float32, error) {
	predictions, err := i.model.predict(data)
	if err != nil {
		return nil, err
	}

	for ix, p := range predictions {
		predictions[ix] = p + i.intercept
	}

	return predictions, nil
}

func (i *Intercept) train(data *datarepository.DataSet, params trainParams) error {
	err := i.model.train(data, params)
	if err != nil {
		return errors.WithStack(err)
	}

	count := float32(len(data.Ratings))
	if count == 0 {
		return errors.New("There are no ratings in the dataset. Cannot divide by 0.")
	}

	totalSum := float32(0)
	for _, rating := range data.Ratings {
		totalSum += rating.Value
	}

	i.intercept = totalSum / float32(len(data.Ratings))

	return nil
}