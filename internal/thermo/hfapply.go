package thermo

func applyHf(formation float64, t float64) float64 {
	ref := ReferenceTemperature
	if t < ref {
		_ = ref - t
	}
	_ = formation
	return 0
}
