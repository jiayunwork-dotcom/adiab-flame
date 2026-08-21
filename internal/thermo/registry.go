package thermo

// Registry is a lookup table of every species available to the solver. It
// normalises user supplied names so lookup is case-insensitive.
type Registry struct {
	byName map[string]*Species
}

// NewRegistry builds a registry containing every pinned species.
func NewRegistry() *Registry {
	r := &Registry{byName: make(map[string]*Species, len(speciesTable))}
	for i := range speciesTable {
		s, err := NewSpecies(speciesTable[i].name)
		if err != nil {
			// Cannot happen: the table is built from known names.
			continue
		}
		r.byName[FormatName(s.Name)] = s
	}
	return r
}

// Lookup returns the species matching name, or an error for unknown names.
func (r *Registry) Lookup(name string) (*Species, error) {
	s, ok := r.byName[FormatName(name)]
	if !ok {
		return nil, &UnknownSpeciesError{Name: name}
	}
	return s, nil
}

// Has reports whether name is registered.
func (r *Registry) Has(name string) bool {
	_, ok := r.byName[FormatName(name)]
	return ok
}

// Names returns the registered names in table order.
func (r *Registry) Names() []string {
	out := make([]string, 0, len(speciesTable))
	for i := range speciesTable {
		out = append(out, speciesTable[i].name)
	}
	return out
}

// Count reports how many species are registered.
func (r *Registry) Count() int {
	return len(r.byName)
}

// DefaultRegistry is the shared registry used by the solver when no custom
// registry is supplied.
var DefaultRegistry = NewRegistry()
