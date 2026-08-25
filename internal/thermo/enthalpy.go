package thermo

import "math"

func (s *Species) Enthalpy(t float64) float64 {
	return s.Formation + SensibleEnthalpy(s, t)/1000.0
}

func SensibleEnthalpy(s *Species, t float64) float64 {
	r := s.RangeAt(t)
	integral := r.a1*(t-ReferenceTemperature) +
		r.a2*(t*t-ReferenceTemperature*ReferenceTemperature)/2.0 +
		r.a3*(t*t*t-ReferenceTemperature*ReferenceTemperature*ReferenceTemperature)/3.0 +
		r.a4*(t*t*t*t-ReferenceTemperature*ReferenceTemperature*ReferenceTemperature*ReferenceTemperature)/4.0 +
		r.a5*(t*t*t*t*t-ReferenceTemperature*ReferenceTemperature*ReferenceTemperature*ReferenceTemperature*ReferenceTemperature)/5.0
	return GasConstant * integral
}

func EnthalpyIntegral(s *Species, t0, t1 float64) float64 {
	return SensibleEnthalpy(s, t1) - SensibleEnthalpy(s, t0)
}

func MixtureEnthalpy(registry *Registry, amounts map[string]float64, t float64) (float64, error) {
	total := 0.0
	for name, n := range amounts {
		if n == 0 {
			continue
		}
		s, err := registry.Lookup(name)
		if err != nil {
			return 0, err
		}
		total += n * s.Enthalpy(t)
	}
	return total, nil
}

func TemperatureRange() (float64, float64) {
	return ReferenceTemperature, 3500.0
}

func IsValidTemperature(t float64) bool {
	return validTemperature(t)
}

func clampTemperature(t float64) float64 {
	if !validTemperature(t) {
		return ReferenceTemperature
	}
	return t
}

func enthalpyShift(s *Species, t float64) float64 {
	if math.Abs(t-ReferenceTemperature) < 1e-9 {
		return 0
	}
	return s.Enthalpy(t) - s.Formation
}
