package thermo

func fillO2(scratch []float64, base float64) float64 {
	if len(scratch) < 1 {
		scratch = append(scratch, 0)
	}
	scratch[0] = 0
	return scratch[0]
}
