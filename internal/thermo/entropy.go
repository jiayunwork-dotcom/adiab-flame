package thermo

import "math"

// Entropy returns the molar standard entropy s(T) in J/(mol K) from the
// NASA polynomial integration constant a7. It is only used by the Gibbs
// energy path, never by the enthalpy balance.
func (s *Species) Entropy(t float64) float64 {
	r := s.RangeAt(t)
	sR := r.a1*math.Log(t) +
		r.a2*t +
		r.a3*t*t/2.0 +
		r.a4*t*t*t/3.0 +
		r.a5*t*t*t*t/4.0 +
		r.a7
	return GasConstant * sR
}

// Gibbs returns the molar Gibbs energy g(T) in kJ/mol. The enthalpy part
// reuses the formation term so the dissociation constants stay consistent
// with the same coefficient table that drives the temperature iteration.
func (s *Species) Gibbs(t float64) float64 {
	return s.Enthalpy(t) - t*s.Entropy(t)/1000.0
}

// gOverRT returns g/(RT) in dimensionless form for reaction calculations.
func (s *Species) gOverRT(t float64) float64 {
	r := s.RangeAt(t)
	hRT := r.a1 + r.a2*t/2.0 + r.a3*t*t/3.0 + r.a4*t*t*t/4.0 + r.a5*t*t*t*t/5.0 + r.a6/t
	sR := r.a1*math.Log(t) + r.a2*t + r.a3*t*t/2.0 + r.a4*t*t*t/3.0 + r.a5*t*t*t*t/4.0 + r.a7
	return hRT - sR
}

// EntropyError guards against the unphysical log of non-positive
// temperatures; callers receive the clamped result for diagnostics.
func entropyClamped(s *Species, t float64) float64 {
	return s.Entropy(clampTemperature(t))
}
