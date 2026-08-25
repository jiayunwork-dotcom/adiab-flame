package thermo

import "fmt"

type Fuel struct {
	Name     string
	Formula  string
	Carbon   float64
	Hydrogen float64
	Species  *Species
}

var fuelTable = []string{"CH4", "C2H4", "C2H6"}

type FuelRegistry struct {
	byName map[string]*Fuel
}

func NewFuelRegistry() *FuelRegistry {
	r := &FuelRegistry{byName: make(map[string]*Fuel, len(fuelTable))}
	for _, name := range fuelTable {
		f := &Fuel{Name: name}
		s, err := NewSpecies(name)
		if err != nil {
			continue
		}
		f.Species = s
		f.Formula = s.Formula
		f.Carbon = s.Atoms["C"]
		f.Hydrogen = s.Atoms["H"]
		r.byName[FormatName(name)] = f
	}
	return r
}

func (r *FuelRegistry) Lookup(name string) (*Fuel, error) {
	f, ok := r.byName[FormatName(name)]
	if !ok {
		return nil, &UnknownFuelError{Name: name}
	}
	return f, nil
}

func (r *FuelRegistry) Has(name string) bool {
	_, ok := r.byName[FormatName(name)]
	return ok
}

func (r *FuelRegistry) Names() []string {
	out := make([]string, 0, len(fuelTable))
	for _, n := range fuelTable {
		out = append(out, n)
	}
	return out
}

func (r *FuelRegistry) Count() int {
	return len(r.byName)
}

func (f *Fuel) Describe() string {
	return fmt.Sprintf("%s (%s)", f.Name, f.Formula)
}

var DefaultFuelRegistry = NewFuelRegistry()
