package thermo

import "math"

// Enthalpy returns the molar enthalpy h(T) in kJ/mol for temperature t,
//
//	h(t) = hf(298.15) + int_{298.15}^{t} cp(tau) d tau
//
// with cp taken from the same pinned NASA polynomial that every other
// thermodynamic function uses. At the reference temperature the enthalpy
// equals the enthalpy of formation exactly.
func (s *Species) Enthalpy(t float64) float64 {
	return s.Formation + SensibleEnthalpy(s, t)/1000.0
}

// SensibleEnthalpy returns int_{T0}^{t} cp d tau in J/mol, the part of the
// enthalpy that is purely thermal. It is the analytic integral of the
// NASA cp polynomial.
func SensibleEnthalpy(s *Species, t float64) float64 {
	r := s.RangeAt(t)
	integral := r.a1*(t-ReferenceTemperature) +
		r.a2*(t*t-ReferenceTemperature*ReferenceTemperature)/2.0 +
		r.a3*(t*t*t-ReferenceTemperature*ReferenceTemperature*ReferenceTemperature)/3.0 +
		r.a4*(t*t*t*t-ReferenceTemperature*ReferenceTemperature*ReferenceTemperature*ReferenceTemperature)/4.0 +
		r.a5*(t*t*t*t*t-ReferenceTemperature*ReferenceTemperature*ReferenceTemperature*ReferenceTemperature*ReferenceTemperature)/5.0
	return GasConstant * integral
}

// EnthalpyIntegral returns h(t) - h(t0) between two temperatures in J/mol,
// recomputed as the difference of two reference integrals.
func EnthalpyIntegral(s *Species, t0, t1 float64) float64 {
	return SensibleEnthalpy(s, t1) - SensibleEnthalpy(s, t0)
}

// MixtureEnthalpy returns the total enthalpy in kJ/mol for a stream given
// as species name to mole amount pairs. The basis is whatever the caller
// uses; results are meaningful as differences between two streams.
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

// TemperatureRange returns a safe bracket for enthalpy monotonicity checks:
// the low edge is the reference temperature, the high edge the polynomial
// validity ceiling.
func TemperatureRange() (float64, float64) {
	return ReferenceTemperature, 3500.0
}

// IsValidTemperature guards all enthalpy evaluations against values the
// polynomial was never fitted for.
func IsValidTemperature(t float64) bool {
	return validTemperature(t)
}

// clampTemperature pulls t into the valid window without panicking the
// caller; NaN inputs fall back to the reference temperature.
func clampTemperature(t float64) float64 {
	if !validTemperature(t) {
		return ReferenceTemperature
	}
	return t
}

// enthalpyShift returns h(t) - hf, i.e. the thermal contribution alone,
// used by diagnostics that must isolate the formation term.
func enthalpyShift(s *Species, t float64) float64 {
	if math.Abs(t-ReferenceTemperature) < 1e-9 {
		return 0
	}
	return s.Enthalpy(t) - s.Formation
}
