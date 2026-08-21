package flame

// ReactantStream builds the inlet composition per mole of fuel: one mole of
// fuel vapour plus the O2 and N2 carried by the actual air supply at the
// given equivalence ratio. The molar basis is defined once here and every
// product stream uses the same basis.
func ReactantStream(cfg *Config, st Stoichiometry) *Stream {
	phi := cfg.EquivalenceRatio
	s := NewStream()
	s.Set(st.FuelName, 1.0)
	s.Set("O2", st.OxygenSupply(phi))
	s.Set("N2", st.NitrogenSupply(phi))
	s.Fuel = st.FuelName
	return s
}

// ProductStream builds the complete-combustion product composition for the
// equivalence ratio. Lean and stoichiometric mixtures fully oxidise the
// fuel; rich mixtures burn only the 1/phi fraction that the supplied oxygen
// can support, leaving unburned fuel vapour in the products. The product
// set is pinned to CO2, H2O, O2, N2 (plus fuel when rich); the optional
// CO/H2 path is applied later by the dissociation refinement.
func ProductStream(cfg *Config, st Stoichiometry) *Stream {
	phi := cfg.EquivalenceRatio
	burned := st.BurnedFraction(phi)

	s := NewStream()
	s.Set("CO2", st.Carbon*burned)
	s.Set("H2O", st.Hydrogen/2.0*burned)
	s.Set("O2", st.ExcessOxygen(phi))
	s.Set("N2", st.NitrogenSupply(phi))
	if burned < 1.0 {
		s.Set(st.FuelName, 1.0-burned)
		s.Fuel = st.FuelName
	}
	return s
}

// productSpecies lists the species the product stream can contain, in
// display order.
var productSpecies = []string{"CO2", "H2O", "O2", "N2", "CO", "H2"}
