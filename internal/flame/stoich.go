package flame

import (
	"fmt"

	"adiab-flame/internal/thermo"
)

type Stoichiometry struct {
	FuelName     string
	Carbon       float64
	Hydrogen     float64
	OxygenNeeded float64
	AirPerFuel   float64
}

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

func (s Stoichiometry) OxygenSupply(phi float64) float64 {
	return s.OxygenNeeded / phi
}

func (s Stoichiometry) NitrogenSupply(phi float64) float64 {
	return thermo.NitrogenForOxygen(s.OxygenSupply(phi))
}

func (s Stoichiometry) AirSupply(phi float64) float64 {
	return thermo.AirForOxygen(s.OxygenSupply(phi))
}

func (s Stoichiometry) ExcessOxygen(phi float64) float64 {
	if phi >= 1.0 {
		return 0
	}
	return s.OxygenNeeded * (1.0/phi - 1.0)
}

func (s Stoichiometry) BurnedFraction(phi float64) float64 {
	if phi <= 1.0 {
		return 1.0
	}
	return 1.0 / phi
}

func (s Stoichiometry) Describe() string {
	return fmt.Sprintf("%s needs %.3f mol O2/mol fuel (air %.3f mol)", s.FuelName, s.OxygenNeeded, s.AirPerFuel)
}
