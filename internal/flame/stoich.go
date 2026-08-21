package flame

import (
	"fmt"

	"adiab-flame/internal/thermo"
)

// Stoichiometry holds the air requirements of one mole of fuel.
type Stoichiometry struct {
	FuelName     string
	Carbon       float64
	Hydrogen     float64
	OxygenNeeded float64 // O2 moles for complete combustion of 1 mol fuel
	AirPerFuel   float64 // air moles for the stoichiometric mixture
}

// ComputeStoichiometry derives the air requirements from the fuel formula.
// The fuel and air molar bases live here once: every later molar stream is
// expressed in "per mole of fuel" and never redefines the air ratio.
func ComputeStoichiometry(fuel *thermo.Fuel) Stoichiometry {
	oxygen := fuel.Carbon + fuel.Hydrogen/4.0
	return Stoichiometry{
		FuelName:     fuel.Name,
		Carbon:       fuel.Carbon,
		Hydrogen:     fuel.Hydrogen,
		OxygenNeeded: oxygen,
		AirPerFuel:   thermo.AirForOxygen(oxygen),
	}
}

// OxygenSupply returns the O2 moles supplied per mole of fuel at the
// equivalence ratio: the stoichiometric demand divided by phi.
func (s Stoichiometry) OxygenSupply(phi float64) float64 {
	return s.OxygenNeeded / phi
}

// NitrogenSupply returns the N2 moles carried with the supplied air.
func (s Stoichiometry) NitrogenSupply(phi float64) float64 {
	return thermo.NitrogenForOxygen(s.OxygenSupply(phi))
}

// AirSupply returns the total air moles supplied per mole of fuel.
func (s Stoichiometry) AirSupply(phi float64) float64 {
	return thermo.AirForOxygen(s.OxygenSupply(phi))
}

// ExcessOxygen returns the O2 moles left over after complete combustion.
// It is zero for rich mixtures.
func (s Stoichiometry) ExcessOxygen(phi float64) float64 {
	if phi >= 1.0 {
		return 0
	}
	return s.OxygenNeeded * (1.0/phi - 1.0)
}

// BurnedFraction is the fraction of fuel that is fully oxidised: 1 for
// lean and stoichiometric mixtures, 1/phi for rich mixtures where the
// oxygen runs out first.
func (s Stoichiometry) BurnedFraction(phi float64) float64 {
	if phi <= 1.0 {
		return 1.0
	}
	return 1.0 / phi
}

// Describe renders the stoichiometry as text for reports.
func (s Stoichiometry) Describe() string {
	return fmt.Sprintf("%s needs %.3f mol O2/mol fuel (air %.3f mol)", s.FuelName, s.OxygenNeeded, s.AirPerFuel)
}
