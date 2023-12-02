package modelevaluation

import (
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestCalculateMae(t *testing.T) {
	testSet := []*datarepository.Rating{
		{Value: 3},
		{Value: 4},
	}
	predictions := []float32{3.5, 1}
	expected := 1.75

	actual, err := calculateMae(testSet, predictions)
	require.NoError(t, err)
	require.Equal(t, true, math.Abs(expected-float64(actual)) <= 0.001)
}

func TestCalculateMaeMissingPrediction(t *testing.T) {
	testSet := []*datarepository.Rating{
		{Value: 3},
		{Value: 4},
	}
	predictions := []float32{3.5}
	_, err := calculateMae(testSet, predictions)
	require.Error(t, err)
}
