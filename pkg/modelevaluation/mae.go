package modelevaluation

import (
	"math"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

func calculateMae(testSet []*datarepository.Rating, predictions []float32) (float32, error) {
	if len(testSet) != len(predictions) {
		return 0, errors.New("Lengths of ratings and predictions are not equal.")
	}

	absError := float64(0)
	for ix, value := range testSet {
		absError = absError + math.Abs(float64(value.Value)-float64(predictions[ix]))
	}
	mae := absError / float64(len(testSet))
	return float32(mae), nil
}
