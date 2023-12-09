package models

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

type DirectorEffect struct {
	model IModel
	directorEffect map[string]float32
}

func (m *DirectorEffect) predict(data *datarepository.DataSet) ([]float32, error) {
	predictions, err := m.model.predict(data)
	if err != nil {
		return nil, err
	}

	for ix, p := range predictions {
		predictions[ix] = p + m.directorEffect[data.Movies[data.Ratings[ix].MovieId].Director]
	}

	return predictions, nil
}

func (m *DirectorEffect) train(data *datarepository.DataSet, params trainParams) error {
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
		factorEffectService.addFactor(data.Movies[rating.MovieId].Director, rating.Value - predictions[ix])
	}

	m.directorEffect = factorEffectService.computeEffects(params.lambda)

	return nil
}