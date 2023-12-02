package datarepository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMovie(t *testing.T) {
	records := [][]string{
		{"1", "292", "5", "838983421", "Outbreak (1995)", "Action|Drama|Sci-Fi|Thriller"},
		{"1", "589", "5", "838983778", "Terminator 2: Judgment Day (1991)", "Action|Sci-Fi"},
		{"2", "539", "3", "868246262", "Sleepless in Seattle (1993)", "Comedy|Drama|Romance"},
		{"3", "540", "2.5", "868246262", "A (1993)", "Romance"},
	}
	expected := []Movie{
		{MovieId: "292", Year: 1995, Title: "Outbreak", Genres: []string{"Action", "Drama", "Sci-Fi", "Thriller"}},
		{MovieId: "589", Year: 1991, Title: "Terminator 2: Judgment Day", Genres: []string{"Action", "Sci-Fi"}},
		{MovieId: "539", Year: 1993, Title: "Sleepless in Seattle", Genres: []string{"Comedy", "Drama", "Romance"}},
		{MovieId: "540", Year: 1993, Title: "A", Genres: []string{"Romance"}},
	}

	for ix, record := range records {
		actual, err := getMovie(record)
		require.NoError(t, err)
		require.Equal(t, expected[ix], *actual)
	}
}

func TestGetRating(t *testing.T) {
	records := [][]string{
		{"1", "292", "5", "838983421", "Outbreak (1995)", "Action|Drama|Sci-Fi|Thriller"},
		{"1", "589", "5", "838983778", "Terminator 2: Judgment Day (1991)", "Action|Sci-Fi"},
		{"2", "539", "3", "868246262", "Sleepless in Seattle (1993)", "Comedy|Drama|Romance"},
		{"3", "540", "2.5", "868246262", "A (1993)", "Romance"},
	}
	expected := []Rating{
		{MovieId: "292", UserId: "1", Value: 5},
		{MovieId: "589", UserId: "1", Value: 5},
		{MovieId: "539", UserId: "2", Value: 3},
		{MovieId: "540", UserId: "3", Value: 2.5},
	}

	for ix, record := range records {
		actual, err := getRating(record)
		require.NoError(t, err)
		require.Equal(t, expected[ix], *actual)
	}
}
