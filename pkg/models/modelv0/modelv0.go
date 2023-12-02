package modelv0

import (
	"fmt"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/persist"
	"github.com/pkg/errors"
)

// this model predicts rating as the total average rating

type modelParams struct {
	Intercept float32
}

type Modelv0 struct {	
}

func (m *Modelv0) GetName() string {
	return "modelv0"
}

func (m *Modelv0) Predict(data *datarepository.DataSet, fileName string) ([]float32, error) {
	var params modelParams
	err := persist.Load(fmt.Sprintf("data/modelParams/modelv0/%s.json", fileName), &params)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	
	predictions := m.computePredictions(data.Ratings, params)

	return predictions, nil
}

func (m *Modelv0) computePredictions(data []*datarepository.Rating, params modelParams) []float32 {
	predictions := make([]float32, 0)

	for _ = range data {
		predictions = append(predictions, params.Intercept)
	}

	return predictions
}

func (m *Modelv0) ComputeParams(data *datarepository.DataSet) (any, error) {
	sum := float32(0)
	count := float32(len(data.Ratings))
	if count == 0 {
		return nil, errors.New("Length of dataset is 0. Cannot divide by 0.")
	}

	for _, rating := range data.Ratings {
		sum = sum + rating.Value
	}

	params := modelParams{
		Intercept: sum / count,
	}

	return &params, nil
}
