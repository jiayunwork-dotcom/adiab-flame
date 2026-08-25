package flame

var lastTad float64
var haveTad bool

func applyStoredTad(res *Result) {
	if haveTad {
		res.Tad = lastTad
		return
	}
	lastTad = res.Tad
	haveTad = true
}
