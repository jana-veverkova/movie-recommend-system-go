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

type FactorObserv struct {
	Factors []string
	Yhat    float32
}

type Model interface {
	ComputeParams(*datarepository.DataSet) (any, error)
	Predict(*datarepository.DataSet, string) ([]float32, error)
	GetName() string
}

func Train(m Model, dataSourceUrl string) error {
	data := datarepository.GetData(dataSourceUrl)

	params, err := m.ComputeParams(data)
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
	data := datarepository.GetData(dataSourceUrl)

	// load parameters
	_, file := path.Split(trainedOnUrl)
	fileName := strings.Split(file, ".")[0]

	// create predictions
	predictions, err := m.Predict(data, fileName)
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

func ComputeFactorEffects(factors []FactorObserv, lambda float32) map[string]float32 {

	factorEffects := make(map[string]float32)
	factorCounts := make(map[string]float32)

	for _, observ := range factors {
		for _, factor := range observ.Factors {
			if factor == "" {
				continue
			}
			factorEffects[factor] = factorEffects[factor] + observ.Yhat
			factorCounts[factor]++
		}
	}

	for key, effect := range factorEffects {
		factorEffects[key] = effect / (factorCounts[key] + lambda)
	}

	return factorEffects
}

func ComputeAverageEffect(effectsmap map[string]float32) (float32, error) {
	if len(effectsmap) == 0 {
		return 0, errors.New("Length of dataset is 0. Cannot divide by 0.")
	}

	sum := float32(0)
	count := float32(0)

	for _, effect := range effectsmap {
		sum = sum + effect
		count++
	}
	
	return sum / count, nil
}

func ComputeIntercept(ratings []*datarepository.Rating) (float32, error) {
	count := float32(len(ratings))
	if count == 0 {
		return 0, errors.New("Length of dataset is 0. Cannot divide by 0.")
	}

	totalSum := float32(0)
	for _, rating := range ratings {
		totalSum += rating.Value
	}
	intercept := totalSum / float32(len(ratings))

	return intercept, nil
}
