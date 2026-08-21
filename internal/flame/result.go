package flame

// Result carries every number a caller needs from a completed calculation:
// the flame temperature, the product composition, the atom residuals and
// the enthalpy balance that closed.
type Result struct {
	Config           Config
	Fuel             string
	Stoichiometry    Stoichiometry
	InletTemperature float64
	EquivalenceRatio float64
	PressureAtm      float64

	Tad              float64
	Reactants        *Stream
	BaseProducts     *Stream
	Products         *Stream
	ReactantEnthalpy float64
	ProductEnthalpy  float64
	Residuals        []AtomResidual
	Iterations       int
	Converged        bool
}

// MoleFractions returns the product mole fractions keyed by species name.
func (r *Result) MoleFractions() map[string]float64 {
	return r.Products.MoleFractions()
}

// MoleFraction returns the product mole fraction of one species.
func (r *Result) MoleFraction(name string) float64 {
	return r.Products.MoleFraction(name)
}

// MaxResidual reports the largest atom-balance residual.
func (r *Result) MaxResidual() float64 {
	return MaxResidual(r.Residuals)
}

// ResidualsNegligible reports whether every atom residual is below the
// pinned threshold.
func (r *Result) ResidualsNegligible() bool {
	return ResidualsNegligible(r.Residuals)
}

// EnthalpyResidual returns h_products(Tad) - h_reactants(Tin), which the
// iteration drives toward zero.
func (r *Result) EnthalpyResidual() float64 {
	return r.ProductEnthalpy - r.ReactantEnthalpy
}

// NitrogenBalance returns the N2 moles entering and leaving, for the
// conservation check.
func (r *Result) NitrogenBalance() (inlet, outlet float64) {
	return r.Reactants.Get("N2"), r.Products.Get("N2")
}

// DissociationUsed reports whether the result went through the optional
// equilibrium refinement.
func (r *Result) DissociationUsed() bool {
	return r.Config.Dissociation
}

// ThermalSpecies returns the species that carry the product enthalpy
// calculation, used for diagnostics.
func (r *Result) ThermalSpecies() []string {
	out := make([]string, 0, len(productSpecies))
	for _, name := range productSpecies {
		if r.Products.SpeciesPresent(name) {
			out = append(out, name)
		}
	}
	return out
}

// ProductAmount returns the molar amount of a product species.
func (r *Result) ProductAmount(name string) float64 {
	return r.Products.Get(name)
}
