package modelevaluation

import (
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestCalculateRmse(t *testing.T) {
	testSet := map[string]datarepository.Rating{
		"1/1" : {UserId: "1", MovieId: "1",	Value: 3,},
		"2/3" : {UserId: "2", MovieId: "3",	Value: 4,},
	}
	predictions := map[string]float32{
		"1/1": 3.5,
		"2/3": 1,
	}
	expected := 2.150581

	actual, err := calculateRmse(testSet, predictions)
	require.NoError(t, err)
	require.Equal(t, true, math.Abs(expected -float64(actual)) <= 0.001)
}

func TestCalculateRmseMissingPrediction(t *testing.T) {
	testSet := map[string]datarepository.Rating{
		"1/1" : {UserId: "1", MovieId: "1",	Value: 3,},
		"2/3" : {UserId: "2", MovieId: "3",	Value: 4,},
	}
	predictions := map[string]float32{
		"1/1": 3.5,
	}
	_, err := calculateRmse(testSet, predictions)
	require.Error(t, err)
}