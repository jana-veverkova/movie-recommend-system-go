package modelv0

import (
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestComputeParams(t *testing.T) {
	var m Modelv0

	ratings := map[string]datarepository.Rating{
		"1": {UserId:  "123", MovieId: "123", Value: 3.5,},
		"2": {UserId:  "1",	MovieId: "2", Value: 4,},
		"3": {UserId:  "2",	MovieId: "3", Value: 5,},
	}
	expected := (3.5 + 4 + 5) / 3

	actualAny, err := m.ComputeParams(ratings)
	require.NoError(t, err)
	actual, ok := actualAny.(*modelParams)
	require.Equal(t, true, ok)
	require.Equal(t, true, math.Abs(expected-float64(actual.Intercept)) <= 0.001)
}

func TestComputeParamsEmptySet(t *testing.T) {
	var m Modelv0
	
	ratings := map[string]datarepository.Rating{}
	_, err := m.ComputeParams(ratings)
	require.Error(t, err)
}

func TestComputePredictions(t *testing.T) {
	var m Modelv0
	
	intercept := float32(0.5)
	data := map[string]datarepository.Rating{
		"1": {UserId:  "123", MovieId: "123", Value: 3.5,},
		"2": {UserId:  "1",	MovieId: "2", Value: 4,},
	}
	expected := map[string]float32{
		"1": 0.5,
		"2": 0.5,
	}

	require.Equal(t, expected, m.computePredictions(data, modelParams{Intercept: intercept}))
}