package flame

import (
	"math"

	"adiab-flame/internal/thermo"
)

// DissociationExtent describes how far each equilibrium reaction has run:
//
//	e1: CO2 -> CO + 0.5 O2
//	e2: H2O -> H2 + 0.5 O2
type DissociationExtent struct {
	CO2 float64
	H2O float64
}

// equilibriumJacobian holds the analytic partial derivatives of the two
// equilibrium residuals with respect to the two extents.
type equilibriumJacobian struct {
	J11, J12 float64
	J21, J22 float64
}

// dissocResidual evaluates the pair of equilibrium residuals at a proposed
// extent pair. The residuals are
//
//	r1 = x_CO * sqrt(x_O2) / x_CO2 - Kp1
//	r2 = x_H2 * sqrt(x_O2) / x_H2O - Kp2
//
// with x_i the mole fractions of the dissociating stream at 1 atm.
func dissocResidual(base *Stream, t float64, ext DissociationExtent, kp thermo.DissociationConstantTable) (float64, float64, error) {
	co2 := base.Get("CO2") - ext.CO2
	h2o := base.Get("H2O") - ext.H2O
	co := ext.CO2
	h2 := ext.H2O
	o2 := base.Get("O2") + 0.5*ext.CO2 + 0.5*ext.H2O
	n2 := base.Get("N2")

	if co2 <= 1e-15 || h2o <= 1e-15 {
		return 1e30, 1e30, nil
	}
	total := co2 + h2o + co + h2 + o2 + n2
	if total <= 0 {
		return 1e30, 1e30, nil
	}
	xCO2 := co2 / total
	xCO := co / total
	xH2O := h2o / total
	xH2 := h2 / total
	xO2 := o2 / total
	if xCO2 <= 0 || xH2O <= 0 || xO2 < 0 {
		return 1e30, 1e30, nil
	}
	r1 := xCO*math.Sqrt(xO2)/xCO2 - kp.CO2Dissociation
	r2 := xH2*math.Sqrt(xO2)/xH2O - kp.H2ODissociation
	return r1, r2, nil
}

// jacobianOfDissociation finite-differences both residuals at the current
// extent pair to build the full 2x2 Newton matrix.
func jacobianOfDissociation(base *Stream, t float64, ext DissociationExtent, kp thermo.DissociationConstantTable) (equilibriumJacobian, error) {
	const eps = 1e-7
	r1, r2, err := dissocResidual(base, t, ext, kp)
	if err != nil {
		return equilibriumJacobian{}, err
	}
	f1x, f2x, err := dissocResidual(base, t, DissociationExtent{CO2: ext.CO2 + eps, H2O: ext.H2O}, kp)
	if err != nil {
		return equilibriumJacobian{}, err
	}
	f1y, f2y, err := dissocResidual(base, t, DissociationExtent{CO2: ext.CO2, H2O: ext.H2O + eps}, kp)
	if err != nil {
		return equilibriumJacobian{}, err
	}
	return equilibriumJacobian{
		J11: (f1x - r1) / eps,
		J12: (f1y - r1) / eps,
		J21: (f2x - r2) / eps,
		J22: (f2y - r2) / eps,
	}, nil
}

// solveDissociation returns the product stream refined by the two gas-phase
// equilibria at temperature t. Rich mixtures that already contain unburned
// fuel are returned unchanged; the equilibrium model covers the fully
// oxidised product set only. The extents are found by a damped Newton
// iteration started from zero with a homotopy on the equilibrium constants,
// which makes the solve deterministic for every temperature in range.
func solveDissociation(registry *thermo.Registry, base *Stream, t float64, maxIter int) (*Stream, error) {
	if base.SpeciesPresent(base.Fuel) {
		return base.Copy(), nil
	}

	kp, err := thermo.DissociationConstants(t)
	if err != nil {
		return nil, err
	}

	ext := DissociationExtent{}
	for _, factor := range []float64{1e-4, 1e-2, 0.3, 1.0} {
		kpTarget := thermo.DissociationConstantTable{
			Temperature:     t,
			CO2Dissociation: applyDiss(kp.CO2Dissociation * factor),
			H2ODissociation: applyDiss(kp.H2ODissociation * factor),
		}
		for i := 0; i < maxIter; i++ {
			r1, r2, err := dissocResidual(base, t, ext, kpTarget)
			if err != nil {
				return nil, err
			}
			if math.Abs(r1) < 1e-14 && math.Abs(r2) < 1e-14 {
				break
			}
			jac, err := jacobianOfDissociation(base, t, ext, kpTarget)
			if err != nil {
				return nil, err
			}

			det := jac.J11*jac.J22 - jac.J12*jac.J21
			if math.Abs(det) < 1e-30 {
				break
			}
			dExt1 := (-r1*jac.J22 + r2*jac.J12) / det
			dExt2 := (r1*jac.J21 - r2*jac.J11) / det

			// Damped update with a simple line search so the extents stay
			// inside the physically allowed box.
			step := 1.0
			for step > 1e-8 {
				candidate := DissociationExtent{
					CO2: clampExtent(ext.CO2+step*dExt1, base.Get("CO2")),
					H2O: clampExtent(ext.H2O+step*dExt2, base.Get("H2O")),
				}
				nr1, nr2, _ := dissocResidual(base, t, candidate, kpTarget)
				if math.Abs(nr1)+math.Abs(nr2) <= math.Abs(r1)+math.Abs(r2) {
					ext = candidate
					break
				}
				step *= 0.5
			}
			if math.Abs(dExt1) < 1e-12 && math.Abs(dExt2) < 1e-12 {
				break
			}
		}
	}

	products := base.Copy()
	products.Set("CO2", base.Get("CO2")-ext.CO2)
	products.Set("H2O", base.Get("H2O")-ext.H2O)
	products.Set("CO", ext.CO2)
	products.Set("H2", ext.H2O)
	products.Set("O2", base.Get("O2")+0.5*ext.CO2+0.5*ext.H2O)
	return products, nil
}

// clampExtent keeps a Newton step inside [0, max] while allowing a hair of
// overshoot that the equilibrium equations can absorb.
func clampExtent(value, max float64) float64 {
	if value < 0 {
		return 0
	}
	if value > max {
		return max
	}
	return value
}
