package models

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

type ActorEffect struct {
	model IModel
	actorEffect map[string]float32
}

func (m *ActorEffect) predict(data *datarepository.DataSet) ([]float32, error) {
	predictions, err := m.model.predict(data)
	if err != nil {
		return nil, err
	}

	for ix, p := range predictions {
		actor1Effect := m.actorEffect[data.Movies[data.Ratings[ix].MovieId].Actor1]
		actor2Effect := m.actorEffect[data.Movies[data.Ratings[ix].MovieId].Actor2]
		predictions[ix] = p + (actor1Effect + actor2Effect)/float32(2)
	}

	return predictions, nil																																									
}

func (m *ActorEffect) train(data *datarepository.DataSet, params trainParams) error {
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
		factorEffectService.addFactor(data.Movies[rating.MovieId].Actor1, rating.Value - predictions[ix])
		factorEffectService.addFactor(data.Movies[rating.MovieId].Actor2, rating.Value - predictions[ix])
	}

	m.actorEffect = factorEffectService.computeEffects(params.lambda)

	return nil
}