package models

import (
	"fmt"
	"math"
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestComputeFactorEffects(t *testing.T) {
	factors := []FactorObserv{
		{Factors: []string{"a"}, Yhat: 0.5},
		{Factors: []string{"a"}, Yhat: 1},
		{Factors: []string{"a"}, Yhat: 2},
		{Factors: []string{"a"}, Yhat: 4},
		{Factors: []string{"b", "a"}, Yhat: 0},
		{Factors: []string{"b"}, Yhat: 2},
		{Factors: []string{"b"}, Yhat: 3},
		{Factors: []string{"b"}, Yhat: 5},
		{Factors: []string{"b"}, Yhat: 3.5},
		{Factors: []string{"c"}, Yhat: 1},
		{Factors: []string{"c"}, Yhat: 3},
	}

	lambda := []float32{0, 10}
	expected := []map[string]float32{
		{"a": 1.5, "b": 2.7, "c": 2},
		{"a": 0.5, "b": 0.9, "c": 0.33333},
	}

	for ix, exp := range expected {
		actual := ComputeFactorEffects(factors, lambda[ix])
		for key, val := range exp {
			require.Equal(t, true, math.Abs(float64(val)-float64(actual[key])) <= 0.001,
				fmt.Sprintf("key: %s, actual factor effect: %.2f, expected: %.2f", key, actual[key], val))
		}
	}
}

func TestAverageEffect(t *testing.T) {
	effectsMap := map[string]float32{
		"a": 1.5,
		"b": 2.7,
		"c": 3,
		"d": 1.1,
	}

	expected := 2.075

	actual, err := ComputeAverageEffect(effectsMap)
	require.NoError(t, err)
	require.Equal(t, true, math.Abs(float64(actual)-float64(expected)) <= 0.001,
		fmt.Sprintf("actual: %.2f, expected: %.2f", actual, expected))
}

func TestComputeIntercept(t *testing.T) {
	ratings := []*datarepository.Rating{
		{Value: 2},
		{Value: 4},
		{Value: 1.5},
	}

	expected := float32(2.5)

	actual, err := ComputeIntercept(ratings)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeInterceptError(t *testing.T) {
	ratings := []*datarepository.Rating{}
	_, err := ComputeIntercept(ratings)
	require.Error(t, err)
}
