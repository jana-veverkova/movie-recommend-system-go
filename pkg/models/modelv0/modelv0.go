package modelv0

import (
	"fmt"
	"path"
	"strings"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/modelevaluation"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/persist"
	"github.com/pkg/errors"
)

type modelParams struct {
	Intercept float32
}

func Train(dataSourceUrl string) error {
	// this model predicts rating as the total average rating

	data, err := datarepository.GetData(dataSourceUrl)
	if err != nil {
		return errors.WithStack(err)
	}

	params, err := computeParams(data.Ratings)
	if err != nil {
		return errors.WithStack(err)
	}

	_, file := path.Split(dataSourceUrl)
	fileName := strings.Split(file, ".")[0]

	err = persist.Save(fmt.Sprintf("data/modelParams/modelv0/%s.json", fileName), params)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Evaluate(trainedOnUrl string, dataSourceUrl string) (*modelevaluation.Summary, error) {
	// load parameters
	_, file := path.Split(trainedOnUrl)
	fileName := strings.Split(file, ".")[0]

	var params modelParams
	err := persist.Load(fmt.Sprintf("data/modelParams/modelv0/%s.json", fileName), &params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// get data for prediction
	data, err := datarepository.GetData(dataSourceUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// create predictions
	predictions := predict(data.Ratings, params)

	// evaluate predictions
	summary, err := modelevaluation.Summarize(data.Ratings, predictions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return summary, nil

}

func predict(data map[string]datarepository.Rating, params modelParams) map[string]float32 {
	// make predictions for data based on params
	predictions := make(map[string]float32)

	for key, _ := range data {
		predictions[key] = params.Intercept
	}

	return predictions
}

func computeParams(ratings map[string]datarepository.Rating) (*modelParams, error) {
	sum := float32(0)
	count := float32(len(ratings))
	if count == 0 {
		return nil, errors.New("Length of dataset is 0. Cannot divide by 0.")
	}

	for _, rating := range ratings {
		sum = sum + rating.Value
	}

	params := modelParams{
		Intercept: sum / count,
	}

	return &params, nil
}
