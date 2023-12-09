package models

import (
	"testing"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/stretchr/testify/require"
)

func TestInterceptTrain(t *testing.T) {
	intercept := Intercept{
		model: &Model{},
	}

	testData := datarepository.DataSet{
		Ratings: []*datarepository.Rating{
			{Value: 1}, {Value: 4}, {Value: 5}, {Value: 6}, {Value: 3}, {Value: 2},
		},
	}

	expectedIntercept := float32(3.5)

	err := intercept.train(&testData, trainParams{lambda: 0})

	require.NoError(t, err)
	require.Equal(t, expectedIntercept, intercept.intercept)
}

func TestInterceptPredict(t *testing.T) {
	intercept := Intercept{
		model: &Model{},
		intercept: 3.5,
	}

	data := datarepository.DataSet{
		Ratings: []*datarepository.Rating{
			{Value: 1}, {Value: 4}, {Value: 5},
		},
	}

	expected := []float32{3.5,3.5,3.5,}
	actual, err := intercept.predict(&data)

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}