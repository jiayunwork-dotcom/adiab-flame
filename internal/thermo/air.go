package thermo

import "fmt"

const AirOxygenFraction = 0.21

const AirNitrogenFraction = 0.79

const AirN2PerO2 = AirNitrogenFraction / AirOxygenFraction

const AirPerO2 = 1.0 / AirOxygenFraction

func AirForOxygen(oxygenMoles float64) float64 {
	return oxygenMoles * AirPerO2
}

func NitrogenForOxygen(oxygenMoles float64) float64 {
	return oxygenMoles * AirN2PerO2
}

func OxygenForAir(airMoles float64) float64 {
	return airMoles * AirOxygenFraction
}

func NitrogenForAir(airMoles float64) float64 {
	return airMoles * AirNitrogenFraction
}

func AirSummary() string {
	return fmt.Sprintf("air = 21%% O2 / 79%% N2 (molar), N2/O2 = %.3f", AirN2PerO2)
}
