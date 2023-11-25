package modelevaluation

import (
	"fmt"
	"math"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/pkg/errors"
)

func calculateRmse(testSet map[string]datarepository.Rating, predictions map[string]float32) (float32, error) {
	sumDiffs := float64(0)
	for key := range testSet {
		if prediction, ok := predictions[key]; ok {
			sumDiffs = sumDiffs + math.Pow(math.Abs(float64(testSet[key].Value) - float64(prediction)), 2)
		} else {
			return 0, errors.New(fmt.Sprintf("Prediction for key %s is missing.", key))
		}
		
	}
	mse := sumDiffs/float64(len(testSet))
	return float32(math.Sqrt(mse)), nil
}