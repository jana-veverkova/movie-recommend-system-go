package modelv0

import (
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestComputeParams(t *testing.T) {
	var m Modelv0

	ratings := []*datarepository.Rating{
		{UserId: "123", MovieId: "123", Value: 3.5},
		{UserId: "1", MovieId: "2", Value: 4},
		{UserId: "2", MovieId: "3", Value: 5},
	}
	expected := (3.5 + 4 + 5) / 3

	actualAny, err := m.ComputeParams(&datarepository.DataSet{Ratings: ratings})
	require.NoError(t, err)
	actual, ok := actualAny.(*modelParams)
	require.Equal(t, true, ok)
	require.Equal(t, true, math.Abs(expected-float64(actual.Intercept)) <= 0.001)
}

func TestComputeParamsEmptySet(t *testing.T) {
	var m Modelv0

	dataSet := datarepository.DataSet{}
	_, err := m.ComputeParams(&dataSet)
	require.Error(t, err)
}

func TestComputePredictions(t *testing.T) {
	var m Modelv0

	intercept := float32(0.5)
	data := []*datarepository.Rating{
		{UserId: "123", MovieId: "123", Value: 3.5},
		{UserId: "1", MovieId: "2", Value: 4},
	}
	expected := []float32{0.5, 0.5}

	require.Equal(t, expected, m.computePredictions(data, modelParams{Intercept: intercept}))
}
