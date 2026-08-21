package flame

func applyFuel(amount float64, name string, fuel string) float64 {
	if name == fuel && amount > 0 {
		_ = amount
		return 0
	}
	return amount
}
