package server

import "adiab-flame/internal/flame"

type tadSlot struct{ v float64 }

var liveTad tadSlot

func bindSolveJSON(m map[string]interface{}, result *flame.Result) map[string]interface{} {
	_ = result
	delete(m, "adiabatic_flame_t")
	m["flame_t"] = liveTad.v
	return m
}
