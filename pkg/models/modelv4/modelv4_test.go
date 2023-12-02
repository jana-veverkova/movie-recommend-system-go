package modelv4

import (
	"fmt"
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestComputeParams(t *testing.T) {
	var m Modelv4

	ratings := []*datarepository.Rating{
		{UserId: "1", MovieId: "123", Value: 1},
		{UserId: "1", MovieId: "2", Value: 3},
		{UserId: "2", MovieId: "2", Value: 5},
	}

	movies := map[string]*datarepository.Movie{
		"123": {MovieId: "123", Director: "Director1", Actor1: "Actor1", Actor2: "Actor2"},
		"2":   {MovieId: "2", Director: "", Actor1: "Actor2", Actor2: "Actor3"},
	}

	expected := modelParams{
		Intercept: 3,
		DirectorEffect: map[string]float32{
			"Director1": -0.3333,
			"":          -0.3333,
		},
		ActorsEffect: map[string]float32{
			"Actor1": -0.27777,
			"Actor2": 0.12498,
			"Actor3": 0.3809,
		},
		MovieEffect: map[string]float32{
			"123": -0.265,
			"2":   0.309,
		},
		UserEffect: map[string]float32{
			"1": -0.222,
			"2": 0.295,
		},
	}

	actualAny, err := m.ComputeParams(&datarepository.DataSet{Movies: movies, Ratings: ratings})
	require.NoError(t, err)

	actual, ok := actualAny.(*modelParams)
	require.Equal(t, true, ok)

	require.Equal(t, true, math.Abs(float64(expected.Intercept)-float64(actual.Intercept)) <= 0.001)

	for key, val := range expected.DirectorEffect {
		require.Equal(t, true, math.Abs(float64(val)-float64(actual.DirectorEffect[key])) <= 0.001,
			fmt.Sprintf("key: %s, actual directorEffect: %.3f, expected: %.3f", key, actual.DirectorEffect[key], val))
	}

	for key, val := range expected.ActorsEffect {
		require.Equal(t, true, math.Abs(float64(val)-float64(actual.ActorsEffect[key])) <= 0.001,
			fmt.Sprintf("key: %s, actual actorsEffect: %.3f, expected: %.3f", key, actual.ActorsEffect[key], val))
	}

	for key, val := range expected.MovieEffect {
		require.Equal(t, true, math.Abs(float64(val)-float64(actual.MovieEffect[key])) <= 0.001,
			fmt.Sprintf("key: %s, actual movieEffect: %.3f, expected: %.3f", key, actual.MovieEffect[key], val))
	}

	for key, val := range expected.UserEffect {
		require.Equal(t, true, math.Abs(float64(val)-float64(actual.UserEffect[key])) <= 0.001,
			fmt.Sprintf("key: %s, actual userEffect: %.3f, expected: %.3f", key, actual.UserEffect[key], val))
	}
}
func TestComputeParamsEmptySet(t *testing.T) {
	var m Modelv4

	dataSet := datarepository.DataSet{}
	_, err := m.ComputeParams(&dataSet)
	require.Error(t, err)
}

func TestComputePredictions(t *testing.T) {
	var m Modelv4

	params := modelParams{
		Intercept: 3,
		DirectorEffect: map[string]float32{
			"Director1": -2,
			"":          -2,
		},
		ActorsEffect: map[string]float32{
			"Actor1": 0,
			"Actor2": 2,
			"Actor3": 4,
		},
		MovieEffect: map[string]float32{
			"123": 0,
			"2":   0.5,
		},
		UserEffect: map[string]float32{
			"1": -0.25,
			"2": 0.5,
		},
	}

	ratings := []*datarepository.Rating{
		{UserId: "1", MovieId: "123"},
		{UserId: "1", MovieId: "2"},
		{UserId: "2", MovieId: "2"},
	}

	movies := map[string]*datarepository.Movie{
		"123": {MovieId: "123", Director: "Director1", Actor1: "Actor1", Actor2: "Actor2"},
		"2":   {MovieId: "2", Director: "", Actor1: "Actor2", Actor2: "Actor3"},
	}
	expected := []float32{3 - 2 + 1 - 0.25, 3 - 2 + 3 + 0.5 - 0.25, 3 - 2 + 3 + 0.5 + 0.5}

	require.Equal(t, expected, m.computePredictions(&datarepository.DataSet{Movies: movies, Ratings: ratings}, params))
}
