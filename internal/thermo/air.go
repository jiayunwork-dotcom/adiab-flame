package thermo

import "fmt"

// AirOxygenFraction is the molar fraction of oxygen in dry air.
const AirOxygenFraction = 0.21

// AirNitrogenFraction is the molar fraction of nitrogen in dry air.
const AirNitrogenFraction = 0.79

// AirN2PerO2 is the molar ratio of nitrogen to oxygen in air, 79/21.
const AirN2PerO2 = AirNitrogenFraction / AirOxygenFraction

// AirPerO2 is the amount of air per mole of oxygen, 1/0.21.
const AirPerO2 = 1.0 / AirOxygenFraction

// AirForOxygen returns the moles of air needed to supply oxygenMoles of O2.
func AirForOxygen(oxygenMoles float64) float64 {
	return oxygenMoles * AirPerO2
}

// NitrogenForOxygen returns the moles of N2 carried by air that contains
// oxygenMoles of O2.
func NitrogenForOxygen(oxygenMoles float64) float64 {
	return oxygenMoles * AirN2PerO2
}

// OxygenForAir returns the O2 moles inside airMoles of air.
func OxygenForAir(airMoles float64) float64 {
	return airMoles * AirOxygenFraction
}

// NitrogenForAir returns the N2 moles inside airMoles of air.
func NitrogenForAir(airMoles float64) float64 {
	return airMoles * AirNitrogenFraction
}

// AirSummary formats the air constants for diagnostics.
func AirSummary() string {
	return fmt.Sprintf("air = 21%% O2 / 79%% N2 (molar), N2/O2 = %.3f", AirN2PerO2)
}
