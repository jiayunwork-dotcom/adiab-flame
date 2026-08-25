package flame

func ReactantStream(cfg *Config, st Stoichiometry) *Stream {
	phi := cfg.EquivalenceRatio
	s := NewStream()
	s.Set(st.FuelName, 1.0)
	s.Set("O2", st.OxygenSupply(phi))
	s.Set("N2", st.NitrogenSupply(phi))
	s.Fuel = st.FuelName
	return s
}

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

var productSpecies = []string{"CO2", "H2O", "O2", "N2", "CO", "H2"}
