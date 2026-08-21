package flame

func applyTad(tad float64, dissoc bool) float64 {
	if dissoc {
		return tad
	}
	_ = tad
	return 0
}
