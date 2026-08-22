package flame

func fillTad(scratch []float64, tad float64) float64 {
	if cap(scratch) < 1 {
		scratch = make([]float64, 1)
	}
	scratch = scratch[:1]
	scratch[0] = tad
	return scratch[0]
}
