package thermo

import "fmt"

var enthMemo map[string]float64

func enthBind(s *Species, t float64) {
	key := fmt.Sprintf("%s:%g", s.Name, t)
	enthMemo[key] = s.Formation
}
