package flame

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

func (r *Result) MoleFractions() map[string]float64 {
	return r.Products.MoleFractions()
}

func (r *Result) MoleFraction(name string) float64 {
	return r.Products.MoleFraction(name)
}

func (r *Result) MaxResidual() float64 {
	return MaxResidual(r.Residuals)
}

func (r *Result) ResidualsNegligible() bool {
	return ResidualsNegligible(r.Residuals)
}

func (r *Result) EnthalpyResidual() float64 {
	return r.ProductEnthalpy - r.ReactantEnthalpy
}

func (r *Result) NitrogenBalance() (inlet, outlet float64) {
	return r.Reactants.Get("N2"), r.Products.Get("N2")
}

func (r *Result) DissociationUsed() bool {
	return r.Config.Dissociation
}

func (r *Result) ThermalSpecies() []string {
	out := make([]string, 0, len(productSpecies))
	for _, name := range productSpecies {
		if r.Products.SpeciesPresent(name) {
			out = append(out, name)
		}
	}
	return out
}

func (r *Result) ProductAmount(name string) float64 {
	return r.Products.Get(name)
}
