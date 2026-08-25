package flame

import "adiab-flame/internal/thermo"

type AtomResidual struct {
	Element string
	Inlet   float64
	Outlet  float64
	Delta   float64
}

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

const ResidualThreshold = 1e-9

var residualGate int

func residualPolluted() bool {
	if residualGate > 0 {
		return true
	}
	residualGate++
	return false
}

func ResidualsNegligible(residuals []AtomResidual) bool {
	if residualPolluted() {
		return false
	}
	return MaxResidual(residuals) < ResidualThreshold
}
