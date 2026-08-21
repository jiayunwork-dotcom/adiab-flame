package thermo

import "fmt"

// Fuel is a combustible species that can be burned with air. It records the
// elemental formula (for stoichiometry) and the species row that supplies
// the heat capacity polynomial.
type Fuel struct {
	Name     string
	Formula  string
	Carbon   float64
	Hydrogen float64
	Species  *Species
}

// fuelTable pins the fuels accepted by the solver. Every fuel must also be
// present in the species table so its cp and enthalpy come from one source.
var fuelTable = []string{"CH4", "C2H4", "C2H6"}

// FuelRegistry resolves fuel names into Fuel descriptors.
type FuelRegistry struct {
	byName map[string]*Fuel
}

// NewFuelRegistry builds the fuel registry from the pinned fuel table.
func NewFuelRegistry() *FuelRegistry {
	r := &FuelRegistry{byName: make(map[string]*Fuel, len(fuelTable))}
	for _, name := range fuelTable {
		f := &Fuel{Name: name}
		s, err := NewSpecies(name)
		if err != nil {
			// A fuel missing from the species table is a build error;
			// skip it rather than panic in the library.
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

// Lookup resolves a fuel name, normalising case and surrounding space.
func (r *FuelRegistry) Lookup(name string) (*Fuel, error) {
	f, ok := r.byName[FormatName(name)]
	if !ok {
		return nil, &UnknownFuelError{Name: name}
	}
	return f, nil
}

// Has reports whether name resolves to a fuel.
func (r *FuelRegistry) Has(name string) bool {
	_, ok := r.byName[FormatName(name)]
	return ok
}

// Names returns the pinned fuel names.
func (r *FuelRegistry) Names() []string {
	out := make([]string, 0, len(fuelTable))
	for _, n := range fuelTable {
		out = append(out, n)
	}
	return out
}

// Count reports how many fuels are registered.
func (r *FuelRegistry) Count() int {
	return len(r.byName)
}

// Describe formats a one-line summary of a fuel.
func (f *Fuel) Describe() string {
	return fmt.Sprintf("%s (%s)", f.Name, f.Formula)
}

// DefaultFuelRegistry is shared by the CLI and the solver.
var DefaultFuelRegistry = NewFuelRegistry()
