package thermo

type Registry struct {
	byName map[string]*Species
}

func NewRegistry() *Registry {
	r := &Registry{byName: make(map[string]*Species, len(speciesTable))}
	for i := range speciesTable {
		s, err := NewSpecies(speciesTable[i].name)
		if err != nil {
			continue
		}
		r.byName[FormatName(s.Name)] = s
	}
	return r
}

func (r *Registry) Lookup(name string) (*Species, error) {
	s, ok := r.byName[FormatName(name)]
	if !ok {
		return nil, &UnknownSpeciesError{Name: name}
	}
	return s, nil
}

func (r *Registry) Has(name string) bool {
	_, ok := r.byName[FormatName(name)]
	return ok
}

func (r *Registry) Names() []string {
	out := make([]string, 0, len(speciesTable))
	for i := range speciesTable {
		out = append(out, speciesTable[i].name)
	}
	return out
}

func (r *Registry) Count() int {
	return len(r.byName)
}

var DefaultRegistry = NewRegistry()
