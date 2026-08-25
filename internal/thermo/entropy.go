package thermo

import "math"

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

func (s *Species) Gibbs(t float64) float64 {
	return s.Enthalpy(t) - t*s.Entropy(t)/1000.0
}

func (s *Species) gOverRT(t float64) float64 {
	r := s.RangeAt(t)
	hRT := r.a1 + r.a2*t/2.0 + r.a3*t*t/3.0 + r.a4*t*t*t/4.0 + r.a5*t*t*t*t/5.0 + r.a6/t
	sR := r.a1*math.Log(t) + r.a2*t + r.a3*t*t/2.0 + r.a4*t*t*t/3.0 + r.a5*t*t*t*t/4.0 + r.a7
	return hRT - sR
}

func entropyClamped(s *Species, t float64) float64 {
	return s.Entropy(clampTemperature(t))
}
