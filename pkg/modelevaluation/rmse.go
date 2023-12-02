package modelevaluation

import (
	"math"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

func calculateRmse(testSet []*datarepository.Rating, predictions []float32) (float32, error) {
	if len(testSet) != len(predictions) {
		return 0, errors.New("Lengths of ratings and predictions are not equal.")
	}
	sumDiffs := float64(0)
	for ix, value := range testSet {
		prediction := predictions[ix]
		sumDiffs = sumDiffs + math.Pow(math.Abs(float64(value.Value) - float64(prediction)), 2)		
	}
	mse := sumDiffs/float64(len(testSet))
	return float32(math.Sqrt(mse)), nil
}