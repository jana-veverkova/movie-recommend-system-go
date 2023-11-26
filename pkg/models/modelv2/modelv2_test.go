package modelv2

import (
	"fmt"
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestComputeParams(t *testing.T) {
	var m Modelv2

	ratings := map[string]datarepository.Rating{
		"1": {UserId: "1", MovieId: "123", Value: 1},
		"2": {UserId: "1", MovieId: "2", Value: 3},
		"3": {UserId: "2", MovieId: "2", Value: 5},
	}
	intercept := 3
	expected := modelParams{
		Intercept: float32(intercept),
		MovieEffect: map[string]float32{
			"123": -2,
			"2":   1,
		},
		UserEffect: map[string]float32{
			"1": -0.5,
			"2": 5 - 3 - 1,
		},
	}

	actualAny, err := m.ComputeParams(ratings)
	require.NoError(t, err)

	actual, ok := actualAny.(*modelParams)
	require.Equal(t, true, ok)

	require.Equal(t, true, math.Abs(float64(expected.Intercept)-float64(actual.Intercept)) <= 0.001)

	for key, val := range expected.MovieEffect {
		require.Equal(t, true, math.Abs(float64(val)-float64(actual.MovieEffect[key])) <= 0.001,
			fmt.Sprintf("key: %s, actual movieEffect: %.2f, expected: %.2f", key, actual.MovieEffect[key], val))
	}

	for key, val := range expected.UserEffect {
		require.Equal(t, true, math.Abs(float64(val)-float64(actual.UserEffect[key])) <= 0.001,
			fmt.Sprintf("key: %s, actual userEffect: %.2f, expected: %.2f", key, actual.UserEffect[key], val))
	}
}
func TestComputeParamsEmptySet(t *testing.T) {
	var m Modelv2

	ratings := map[string]datarepository.Rating{}
	_, err := m.ComputeParams(ratings)
	require.Error(t, err)
}

func TestComputePredictions(t *testing.T) {
	var m Modelv2

	params := modelParams{
		Intercept: 0.5,
		MovieEffect: map[string]float32{
			"123": -1,
			"2":   3,
		},
		UserEffect: map[string]float32{
			"1": 1,
			"2": -1.5,
		},
	}
	data := map[string]datarepository.Rating{
		"1": {UserId: "1", MovieId: "123", Value: 3.5},
		"2": {UserId: "1", MovieId: "2", Value: 4},
		"3": {UserId: "2", MovieId: "2", Value: 5},
	}
	expected := map[string]float32{
		"1": 0.5 - 1 + 1,
		"2": 0.5 + 3 + 1,
		"3": 0.5 + 3 - 1.5,
	}

	require.Equal(t, expected, m.computePredictions(data, params))
}
