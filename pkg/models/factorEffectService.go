package models

type FactorEffectService struct {
	factors []string
	yHat    []float32
}

func (s *FactorEffectService) addFactor(factor string, yHat float32) {
	s.factors = append(s.factors, factor)
	s.yHat = append(s.yHat, yHat)
}

func (s *FactorEffectService) computeEffects(lambda float32) map[string]float32 {
	factorEffects := make(map[string]float32)
	factorCounts := make(map[string]float32)
	totalEffect := float32(0)

	for ix, factor := range s.factors {
		if factor == "" {
			continue
		}
		factorEffects[factor] = factorEffects[factor] + s.yHat[ix]
		factorCounts[factor]++
	}

	for key, effect := range factorEffects {
		factorEffects[key] = effect / (factorCounts[key] + lambda)
		totalEffect = totalEffect + effect / (factorCounts[key] + lambda)
	}

	factorEffects[""] = totalEffect/float32(len(factorEffects))

	return factorEffects
}
