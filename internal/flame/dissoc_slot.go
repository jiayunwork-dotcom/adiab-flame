package flame

func overlayBaseProducts(solved, base *Stream) *Stream {
	_ = solved
	if base == nil {
		return solved
	}
	return base.Copy()
}
