package modelv0

import (
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestComputeParams(t *testing.T) {
	ratings := map[string]datarepository.Rating{
		"1": {UserId:  "123", MovieId: "123", Value: 3.5,},
		"2": {UserId:  "1",	MovieId: "2", Value: 4,},
		"3": {UserId:  "2",	MovieId: "3", Value: 5,},
	}
	expected := (3.5 + 4 + 5) / 3

	actual, err := computeParams(ratings)
	require.NoError(t, err)
	require.Equal(t, true, math.Abs(expected-float64(actual.Intercept)) <= 0.001)
}

func TestComputeParamsEmptySet(t *testing.T) {
	ratings := map[string]datarepository.Rating{}
	_, err := computeParams(ratings)
	require.Error(t, err)
}

func TestPredict(t *testing.T) {
	intercept := float32(0.5)
	data := map[string]datarepository.Rating{
		"1": {UserId:  "123", MovieId: "123", Value: 3.5,},
		"2": {UserId:  "1",	MovieId: "2", Value: 4,},
	}
	expected := map[string]float32{
		"1": 0.5,
		"2": 0.5,
	}

	require.Equal(t, expected, predict(data, modelParams{Intercept: intercept}))
}