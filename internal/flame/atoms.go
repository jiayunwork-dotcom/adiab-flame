package flame

import "adiab-flame/internal/thermo"

// AtomResidual is the difference between atoms entering and leaving the
// reactor, on a per-mole-of-fuel basis. A zero residual means the element
// balance closes exactly.
type AtomResidual struct {
	Element string
	Inlet   float64
	Outlet  float64
	Delta   float64
}

// AtomResiduals computes the C/H/O/N balance between the reactant and
// product streams.
func AtomResiduals(registry *thermo.Registry, reactants, products *Stream) ([]AtomResidual, error) {
	inlet, err := reactants.ElementTotals(registry)
	if err != nil {
		return nil, err
	}
	outlet, err := products.ElementTotals(registry)
	if err != nil {
		return nil, err
	}
	out := make([]AtomResidual, 0, len(thermo.Elements))
	for _, el := range thermo.Elements {
		out = append(out, AtomResidual{
			Element: el,
			Inlet:   inlet[el],
			Outlet:  outlet[el],
			Delta:   inlet[el] - outlet[el],
		})
	}
	return out, nil
}

// MaxResidual returns the largest absolute residual magnitude.
func MaxResidual(residuals []AtomResidual) float64 {
	max := 0.0
	for _, r := range residuals {
		if d := abs(r.Delta); d > max {
			max = d
		}
	}
	return max
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// ResidualThreshold is the pinned bound under which residuals are
// considered negligible.
const ResidualThreshold = 1e-9

// ResidualsNegligible reports whether every residual stays below the
// pinned threshold.
func ResidualsNegligible(residuals []AtomResidual) bool {
	return MaxResidual(residuals) < ResidualThreshold
}
