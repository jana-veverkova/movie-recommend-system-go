package modelv2B

import (
	"fmt"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/persist"
	"github.com/pkg/errors"
)

type modelParams struct {
	Intercept               float32
	MovieEffect, UserEffect map[string]float32
}

type Modelv2B struct {
}

func (m *Modelv2B) GetName() string {
	return "modelv2B"
}

func (m *Modelv2B) Predict(data *datarepository.DataSet, fileName string) ([]float32, error) {
	var params modelParams
	err := persist.Load(fmt.Sprintf("data/modelParams/modelv2B/%s.json", fileName), &params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	predictions := m.computePredictions(data.Ratings, params)

	return predictions, nil
}

func (m *Modelv2B) computePredictions(data []*datarepository.Rating, params modelParams) []float32 {
	predictions := make([]float32, 0)

	for _, val := range data {
		predictions = append(predictions, params.Intercept+params.MovieEffect[val.MovieId]+params.UserEffect[val.UserId])
	}

	return predictions
}

func (m *Modelv2B) ComputeParams(data *datarepository.DataSet) (any, error) {
	lambdaMovies := float32(100)
	lambdaUsers := float32(100)

	totalSum := float32(0)
	count := float32(len(data.Ratings))
	if count == 0 {
		return nil, errors.New("Length of dataset is 0. Cannot divide by 0.")
	}

	// compute intercept
	for _, rating := range data.Ratings {
		totalSum += rating.Value
	}

	intercept := totalSum / count

	// compute movies effects
	// compute movies effects
	movieFactors := make([]models.FactorObserv, 0)
	for _, rating := range data.Ratings {
		movieFactors = append(movieFactors, models.FactorObserv{
			Factors: []string{rating.MovieId},
			Yhat:    rating.Value - intercept,
		})
	}
	moviesEffects := models.ComputeFactorEffects(movieFactors, lambdaMovies)
	avgEffect, err := models.ComputeAverageEffect(moviesEffects)
	if err != nil {
		return nil, err
	}
	moviesEffects[""] = avgEffect

	// compute user effect
	userFactors := make([]models.FactorObserv, 0)
	for _, rating := range data.Ratings {
		userFactors = append(userFactors, models.FactorObserv{
			Factors: []string{rating.UserId},
			Yhat:    rating.Value - intercept - moviesEffects[rating.MovieId],
		})
	}

	userEffects := models.ComputeFactorEffects(userFactors, float32(lambdaUsers))
	userEffects[""], err = models.ComputeAverageEffect(userEffects)
	if err != nil {
		return nil, err
	}

	params := modelParams{
		Intercept:   totalSum / count,
		MovieEffect: moviesEffects,
		UserEffect:  userEffects,
	}

	return &params, nil
}
