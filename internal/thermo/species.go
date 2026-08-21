package thermo

import (
	"fmt"
	"math"
	"strings"
)

// GasConstant is the molar gas constant in J/(mol K).
const GasConstant = 8.31446261815324

// ReferenceTemperature is the standard state temperature used for the
// enthalpy of formation, in kelvin.
const ReferenceTemperature = 298.15

// Element names handled by the atom balance. The order is fixed so that
// residual vectors print in a stable order.
var Elements = []string{"C", "H", "O", "N"}

// Species describes one chemical species: its identity, elemental formula
// and the pinned thermodynamic table entry that produces its properties.
type Species struct {
	Name      string
	Formula   string
	Formation float64 // standard enthalpy of formation at 298.15 K, kJ/mol
	Atoms     map[string]float64
	entry     *nasaEntry
}

// NewSpecies builds a Species view from a name. It returns an error when the
// name is not present in the pinned table.
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

// MustSpecies is NewSpecies with a panic on error, for use only with
// hard-coded names that are known to exist in the table.
func MustSpecies(name string) *Species {
	s, err := NewSpecies(name)
	if err != nil {
		panic(err)
	}
	return s
}

// countAtoms parses a formula like "C2H6" into an element map. Element
// symbols are one or two letters followed by an optional count.
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

// RangeAt picks the NASA polynomial range active at temperature t.
func (s *Species) RangeAt(t float64) *nasaRange {
	if t < s.entry.mid {
		return &s.entry.low
	}
	return &s.entry.high
}

// Cappacity returns the molar heat capacity at constant pressure in
// J/(mol K) from the same pinned polynomial used for the enthalpy integral.
func (s *Species) Cappacity(t float64) float64 {
	r := s.RangeAt(t)
	return GasConstant * (r.a1 + r.a2*t + r.a3*t*t + r.a4*t*t*t + r.a5*t*t*t*t)
}

// CpOverR returns cp/R, the polynomial value without the gas constant.
func (s *Species) CpOverR(t float64) float64 {
	r := s.RangeAt(t)
	return r.a1 + r.a2*t + r.a3*t*t + r.a4*t*t*t + r.a5*t*t*t*t
}

// validTemperature reports whether t lies inside the polynomial validity
// window shared by both ranges (200 K to 3500 K covers every table entry).
func validTemperature(t float64) bool {
	return !math.IsNaN(t) && !math.IsInf(t, 0) && t >= 200.0 && t <= 3500.0
}

// FormatName normalises user supplied species names so that "co2", "CO2"
// and " Co2 " all resolve to the same table row.
func FormatName(name string) string {
	return strings.ToUpper(strings.TrimSpace(name))
}
