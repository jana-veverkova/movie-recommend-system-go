package datarepository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatData(t *testing.T) {
	records := [][]string{
		{"1", "292", "5", "838983421", "Outbreak (1995)", "Action|Drama|Sci-Fi|Thriller"},
		{"1", "589", "5", "838983778", "Terminator 2: Judgment Day (1991)", "Action|Sci-Fi"},
		{"2", "539", "3", "868246262", "Sleepless in Seattle (1993)", "Comedy|Drama|Romance"},
		{"2", "539", "3", "868246262", "Sleepless in Seattle (1993)", "Comedy|Drama|Romance"},
		{"3", "540", "2.5", "868246262", "A (1993)", "Romance"},
	}
	expected := dataSet{
		movies: map[string]movie{
			"292": {
				movieId: "292",
				year: 1995,
				title: "Outbreak",
				genres: []string{"Action", "Drama", "Sci-Fi", "Thriller"},
			},
			"589": {
				movieId: "589",
				year: 1991,
				title: "Terminator 2: Judgment Day",
				genres: []string{"Action", "Sci-Fi"},
			},
			"539": {
				movieId: "539",
				year: 1993,
				title: "Sleepless in Seattle",
				genres: []string{"Comedy","Drama","Romance"},
			},
			"540": {
				movieId: "540",
				year: 1993,
				title: "A",
				genres: []string{"Romance"},
			},
		},
		ratings: map[string]rating{
			"1/292": {
				movieId: "292",
				userId: "1",
				rating: 5,
			},
			"1/589": {
				movieId: "589",
				userId: "1",
				rating: 5,
			},
			"2/539": {
				movieId: "539",
				userId: "2",
				rating: 3,
			},
			"3/540": {
				movieId: "540",
				userId: "3",
				rating: 2.5,
			},
		},
	}

	actual, err := formatData(records)
	require.NoError(t, err)
	require.Equal(t, &expected, actual)
}