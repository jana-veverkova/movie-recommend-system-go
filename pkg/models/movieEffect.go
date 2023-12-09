package models

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

type MovieEffect struct {
	model IModel
	movieEffect map[string]float32
}

func (m *MovieEffect) predict(data *datarepository.DataSet) ([]float32, error) {
	predictions, err := m.model.predict(data)
	if err != nil {
		return nil, err
	}

	for ix, p := range predictions {
		predictions[ix] = p + m.movieEffect[data.Ratings[ix].MovieId]
	}

	return predictions, nil
}

func (m *MovieEffect) train(data *datarepository.DataSet, params trainParams) error {
	err := m.model.train(data, params)
	if err != nil {
		return errors.WithStack(err)
	}
	
	predictions, err := m.model.predict(data)
	if err != nil {
		return errors.WithStack(err)
	}

	factorEffectService := FactorEffectService{}

	for ix, rating := range data.Ratings {
		factorEffectService.addFactor(rating.MovieId, rating.Value - predictions[ix])
	}

	m.movieEffect = factorEffectService.computeEffects(params.lambda)

	return nil
}