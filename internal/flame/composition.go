package flame

import (
	"fmt"
	"sort"

	"adiab-flame/internal/thermo"
)

// Stream is a molar composition on a per-mole-of-fuel basis. Amounts are
// stored by species name so reactant and product streams share one type.
type Stream struct {
	// Species maps species name to molar amount.
	Species map[string]float64
	// Fuel records the fuel species name when unburned fuel is present.
	Fuel string
}

// NewStream returns an empty stream.
func NewStream() *Stream {
	return &Stream{Species: make(map[string]float64)}
}

// Set assigns a molar amount to a species.
func (s *Stream) Set(name string, moles float64) {
	if moles < 0 {
		moles = 0
	}
	s.Species[name] = moles
}

// Get returns the molar amount of a species (zero if absent).
func (s *Stream) Get(name string) float64 {
	return s.Species[name]
}

// Add adds moles to a species.
func (s *Stream) Add(name string, moles float64) {
	s.Species[name] += moles
}

// Remove deletes a species entirely.
func (s *Stream) Remove(name string) {
	delete(s.Species, name)
}

// TotalMoles sums every species amount.
func (s *Stream) TotalMoles() float64 {
	total := 0.0
	for _, n := range s.Species {
		total += n
	}
	return total
}

// MoleFraction returns the fraction of a species in the stream.
func (s *Stream) MoleFraction(name string) float64 {
	total := s.TotalMoles()
	if total <= 0 {
		return 0
	}
	return s.Species[name] / total
}

// MoleFractions returns a copy of every species fraction keyed by name.
func (s *Stream) MoleFractions() map[string]float64 {
	total := s.TotalMoles()
	out := make(map[string]float64, len(s.Species))
	for name, n := range s.Species {
		if total > 0 {
			out[name] = n / total
		} else {
			out[name] = 0
		}
	}
	return out
}

// ElementTotals sums the atoms of each element carried by the stream.
func (s *Stream) ElementTotals(registry *thermo.Registry) (map[string]float64, error) {
	totals := make(map[string]float64, len(thermo.Elements))
	for _, el := range thermo.Elements {
		totals[el] = 0
	}
	for name, moles := range s.Species {
		if moles == 0 {
			continue
		}
		spec, err := registry.Lookup(name)
		if err != nil {
			return nil, err
		}
		for _, el := range thermo.Elements {
			totals[el] += moles * spec.Atoms[el]
		}
	}
	return totals, nil
}

// Enthalpy returns the total enthalpy of the stream in kJ per mole of fuel.
func (s *Stream) Enthalpy(registry *thermo.Registry, t float64) (float64, error) {
	total := 0.0
	for name, moles := range s.Species {
		if moles == 0 {
			continue
		}
		spec, err := registry.Lookup(name)
		if err != nil {
			return 0, err
		}
		total += moles * spec.Enthalpy(t)
	}
	return total, nil
}

// SpeciesPresent reports whether a species has a positive amount.
func (s *Stream) SpeciesPresent(name string) bool {
	return s.Species[name] > 0
}

// Describe returns a stable textual rendering: sorted species, one per line.
func (s *Stream) Describe() string {
	names := make([]string, 0, len(s.Species))
	for name := range s.Species {
		names = append(names, name)
	}
	sort.Strings(names)
	out := ""
	for _, name := range names {
		out += fmt.Sprintf("%s %.6g\n", name, s.Species[name])
	}
	return out
}

// Copy deep-copies the stream.
func (s *Stream) Copy() *Stream {
	out := NewStream()
	for name, moles := range s.Species {
		out.Species[name] = moles
	}
	out.Fuel = s.Fuel
	return out
}
