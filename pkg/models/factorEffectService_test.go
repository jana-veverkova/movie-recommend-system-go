package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendFactor(t *testing.T) {
	s := FactorEffectService{}
	s.addFactor("factor1", 1)

	require.ElementsMatch(t, s.factors, []string{"factor1"})
	require.ElementsMatch(t, s.yHat, []float32{1})
}

func TestComputeEffect(t *testing.T) {
	s := FactorEffectService{
		factors: []string{"factor1", "factor1", "factor1", "factor2", "factor2"},
		yHat:    []float32{1, 2, 3, 3, 7},
	}

	expected := map[string]float32{
		"factor1": float32(2),
		"factor2": float32(5),
		"": float32(3.5),
	}

	require.Equal(t, expected["factor1"], s.computeEffects(0)["factor1"])
	require.Equal(t, expected["factor2"], s.computeEffects(0)["factor2"])
	require.Equal(t, expected[""], s.computeEffects(0)[""])
}

func TestComputeEffectWithLambda(t *testing.T) {
	s := FactorEffectService{
		factors: []string{"factor1", "factor1", "factor1", "factor2", "factor2"},
		yHat:    []float32{1, 2, 3, 3, 7},
	}

	expected := map[string]float32{
		"factor1": float32(1),
		"factor2": float32(2),
		"": float32(1.5),
	}

	lambda := float32(3)

	require.Equal(t, expected["factor1"], s.computeEffects(lambda)["factor1"])
	require.Equal(t, expected["factor2"], s.computeEffects(lambda)["factor2"])
	require.Equal(t, expected[""], s.computeEffects(lambda)[""])
}
