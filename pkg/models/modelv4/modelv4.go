package modelv4

import (
	"fmt"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/persist"
	"github.com/pkg/errors"
)

type modelParams struct {
	Intercept                                             float32
	DirectorEffect, MovieEffect, ActorsEffect, UserEffect map[string]float32
}

type Modelv4 struct {
}

func (m *Modelv4) GetName() string {
	return "modelv4"
}

func (m *Modelv4) Predict(data *datarepository.DataSet, fileName string) ([]float32, error) {
	var params modelParams
	err := persist.Load(fmt.Sprintf("data/modelParams/modelv4/%s.json", fileName), &params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	predictions := m.computePredictions(data, params)

	return predictions, nil
}

func (m *Modelv4) computePredictions(data *datarepository.DataSet, params modelParams) []float32 {
	predictions := make([]float32, 0)

	for _, r := range data.Ratings {
		d := params.DirectorEffect[data.Movies[r.MovieId].Director]
		a1 := params.ActorsEffect[data.Movies[r.MovieId].Actor1]
		a2 := params.ActorsEffect[data.Movies[r.MovieId].Actor2]
		m := params.MovieEffect[r.MovieId]
		u := params.UserEffect[r.UserId]
		fmt.Println(d, a1, a2, m, u)
		predictions = append(predictions,
			params.Intercept+d+(a1+a2)/2+m+u)
	}

	return predictions
}

func (m *Modelv4) ComputeParams(data *datarepository.DataSet) (any, error) {
	lambda := float32(5)

	intercept, err := models.ComputeIntercept(data.Ratings)
	if err != nil {
		return nil, err
	}

	// compute director effects
	directorFactors := make([]models.FactorObserv, 0)
	for _, rating := range data.Ratings {
		directorFactors = append(directorFactors, models.FactorObserv{
			Factors: []string{data.Movies[rating.MovieId].Director},
			Yhat:    rating.Value - intercept,
		})
	}

	directorEffects := models.ComputeFactorEffects(directorFactors, lambda)
	directorEffects[""], err = models.ComputeAverageEffect(directorEffects)
	if err != nil {
		return nil, err
	}

	// compute actors effect
	actorsFactors := make([]models.FactorObserv, 0)
	for _, rating := range data.Ratings {
		actorsFactors = append(actorsFactors, models.FactorObserv{
			Factors: []string{data.Movies[rating.MovieId].Actor1, data.Movies[rating.MovieId].Actor2},
			Yhat:    rating.Value - intercept - directorEffects[data.Movies[rating.MovieId].Director],
		})
	}

	actorsEffects := models.ComputeFactorEffects(actorsFactors, lambda)
	actorsEffects[""], err = models.ComputeAverageEffect(actorsEffects)
	if err != nil {
		return nil, err
	}

	// compute movies effects
	movieFactors := make([]models.FactorObserv, 0)
	for _, rating := range data.Ratings {
		movieFactors = append(movieFactors, models.FactorObserv{
			Factors: []string{rating.MovieId},
			Yhat:    rating.Value - intercept - directorEffects[data.Movies[rating.MovieId].Director] - (actorsEffects[data.Movies[rating.MovieId].Actor1]+actorsEffects[data.Movies[rating.MovieId].Actor2])/2,
		})
	}
	moviesEffects := models.ComputeFactorEffects(movieFactors, lambda)
	moviesEffects[""], err = models.ComputeAverageEffect(moviesEffects)
	if err != nil {
		return nil, err
	}

	// compute user effect
	userFactors := make([]models.FactorObserv, 0)
	for _, rating := range data.Ratings {
		userFactors = append(userFactors, models.FactorObserv{
			Factors: []string{rating.UserId},
			Yhat:    rating.Value - intercept - moviesEffects[rating.MovieId] - directorEffects[data.Movies[rating.MovieId].Director] - (actorsEffects[data.Movies[rating.MovieId].Actor1]+actorsEffects[data.Movies[rating.MovieId].Actor2])/2,
		})
	}

	userEffects := models.ComputeFactorEffects(userFactors, lambda)
	userEffects[""], err = models.ComputeAverageEffect(userEffects)
	if err != nil {
		return nil, err
	}

	params := modelParams{
		Intercept:      intercept,
		DirectorEffect: directorEffects,
		MovieEffect:    moviesEffects,
		ActorsEffect:   actorsEffects,
		UserEffect:     userEffects,
	}

	return &params, nil
}
