package modelv2

import (
	"fmt"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/persist"
	"github.com/pkg/errors"
)

type modelParams struct {
	Intercept               float32
	MovieEffect, UserEffect map[string]float32
}

type counts struct {
	sum, count float32
}

type Modelv2 struct {
}

func (m *Modelv2) GetName() string {
	return "modelv2"
}

func (m *Modelv2) Predict(data *datarepository.DataSet, fileName string) ([]float32, error) {
	var params modelParams
	err := persist.Load(fmt.Sprintf("data/modelParams/modelv2/%s.json", fileName), &params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	predictions := m.computePredictions(data.Ratings, params)

	return predictions, nil
}

func (m *Modelv2) computePredictions(data []*datarepository.Rating, params modelParams) []float32 {
	predictions := make([]float32, 0)

	for _, val := range data {
		predictions = append(predictions, params.Intercept + params.MovieEffect[val.MovieId] + params.UserEffect[val.UserId])
	}

	return predictions
}

func (m *Modelv2) ComputeParams(data *datarepository.DataSet) (any, error) {
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
	moviesCounts := make(map[string]counts)
	for _, rating := range data.Ratings {
		yComma := rating.Value - intercept
		moviesCounts[rating.MovieId] = counts{moviesCounts[rating.MovieId].sum + yComma, moviesCounts[rating.MovieId].count + 1}
	}

	moviesEffects := make(map[string]float32)
	for key, effect := range moviesCounts {
		moviesEffects[key] = effect.sum / effect.count
	}

	// compute user effect
	usersCounts := make(map[string]counts)
	for _, rating := range data.Ratings {
		yComma := rating.Value - intercept - float32(moviesEffects[rating.MovieId])
		usersCounts[rating.UserId] = counts{usersCounts[rating.UserId].sum + yComma, usersCounts[rating.UserId].count + 1}
	}

	usersEffects := make(map[string]float32)
	for key, effect := range usersCounts {
		usersEffects[key] = effect.sum / effect.count
	}

	params := modelParams{
		Intercept:   totalSum / count,
		MovieEffect: moviesEffects,
		UserEffect:  usersEffects,
	}

	return &params, nil
}
