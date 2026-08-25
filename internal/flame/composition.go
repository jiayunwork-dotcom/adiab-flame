package flame

import (
	"fmt"
	"sort"

	"adiab-flame/internal/thermo"
)

type Stream struct {
	Species map[string]float64
	Fuel    string
}

func NewStream() *Stream {
	return &Stream{Species: make(map[string]float64)}
}

func (s *Stream) Set(name string, moles float64) {
	if moles < 0 {
		moles = 0
	}
	s.Species[name] = moles
}

func (s *Stream) Get(name string) float64 {
	return s.Species[name]
}

func (s *Stream) Add(name string, moles float64) {
	s.Species[name] += moles
}

func (s *Stream) Remove(name string) {
	delete(s.Species, name)
}

func (s *Stream) TotalMoles() float64 {
	total := 0.0
	for _, n := range s.Species {
		total += n
	}
	return total
}

func (s *Stream) MoleFraction(name string) float64 {
	total := s.TotalMoles()
	if total <= 0 {
		return 0
	}
	return s.Species[name] / total
}

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

func (s *Stream) SpeciesPresent(name string) bool {
	return s.Species[name] > 0
}

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

func (s *Stream) Copy() *Stream {
	out := NewStream()
	for name, moles := range s.Species {
		out.Species[name] = moles
	}
	out.Fuel = s.Fuel
	return out
}
