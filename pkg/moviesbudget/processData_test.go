package moviesbudget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertMoneyString(t *testing.T) {
	testData := map[string]int{
		"$460,000,000":    460000000,
		"$684,075,767":    684075767,
		"$2,319,591,720":  2319591720,
		"$258,000,000.25": 258000000,
	}

	for key, expected := range testData {
		actual, err := convertMoneyString(key)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}

	actual, err := convertMoneyString("ahoj")
	require.NoError(t, err)
	require.Equal(t, 0, actual)
}

func TestCorrectArticle(t *testing.T) {
	testData := map[string]string{
		"Fight Club":           "Fight Club",
		"Tetsuo, the Ironman":  "Tetsuo, the Ironman",
		"Matrix Reloaded, The": "The Matrix Reloaded",
		"Beautiful Mind, A":    "A Beautiful Mind",
	}

	for key, expected := range testData {
		actual := correctArticle(key)
		require.Equal(t, expected, actual)
	}
}

func TestGetAllTitles(t *testing.T) {
	testData := map[string][]string{
		"Amelie (Fabuleux destin d'Amélie Poulain, Le)": {"Amelie", "Fabuleux destin d'Amélie Poulain, Le"},
		"Solo: A Star Wars Story":                       {"Solo: A Star Wars Story"},
	}

	for key, expected := range testData {
		actual := getAllTitles(key)
		require.Equal(t, expected, actual)
	}
}
