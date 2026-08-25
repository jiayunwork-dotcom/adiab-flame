package thermo

import (
	"fmt"
	"math"
	"strings"
)

const GasConstant = 8.31446261815324

const ReferenceTemperature = 298.15

var Elements = []string{"C", "H", "O", "N"}

type Species struct {
	Name      string
	Formula   string
	Formation float64
	Atoms     map[string]float64
	entry     *nasaEntry
}

func NewSpecies(name string) (*Species, error) {
	idx, ok := tableIndex[name]
	if !ok {
		return nil, fmt.Errorf("unknown species %q", name)
	}
	e := &speciesTable[idx]
	return &Species{
		Name:      e.name,
		Formula:   e.formula,
		Formation: e.dhf,
		Atoms:     countAtoms(e.formula),
		entry:     e,
	}, nil
}

func MustSpecies(name string) *Species {
	s, err := NewSpecies(name)
	if err != nil {
		panic(err)
	}
	return s
}

func countAtoms(formula string) map[string]float64 {
	atoms := make(map[string]float64)
	i := 0
	for i < len(formula) {
		start := i
		i++
		for i < len(formula) && formula[i] >= 'a' && formula[i] <= 'z' {
			i++
		}
		sym := formula[start:i]
		digits := ""
		for i < len(formula) && formula[i] >= '0' && formula[i] <= '9' {
			digits += string(formula[i])
			i++
		}
		n := 1.0
		if digits != "" {
			n = float64(parseUint(digits))
		}
		atoms[sym] += n
	}
	return atoms
}

func parseUint(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func (s *Species) RangeAt(t float64) *nasaRange {
	if t < s.entry.mid {
		return &s.entry.low
	}
	return &s.entry.high
}

func (s *Species) Cappacity(t float64) float64 {
	r := s.RangeAt(t)
	return GasConstant * (r.a1 + r.a2*t + r.a3*t*t + r.a4*t*t*t + r.a5*t*t*t*t)
}

func (s *Species) CpOverR(t float64) float64 {
	r := s.RangeAt(t)
	return r.a1 + r.a2*t + r.a3*t*t + r.a4*t*t*t + r.a5*t*t*t*t
}

func validTemperature(t float64) bool {
	return !math.IsNaN(t) && !math.IsInf(t, 0) && t >= 200.0 && t <= 3500.0
}

func FormatName(name string) string {
	return strings.ToUpper(strings.TrimSpace(name))
}
