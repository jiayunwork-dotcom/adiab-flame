package thermo

import "fmt"

// UnknownSpeciesError reports a species name that is not in the pinned
// thermodynamic table.
type UnknownSpeciesError struct {
	Name string
}

func (e *UnknownSpeciesError) Error() string {
	return fmt.Sprintf("unknown species %q (known: CH4, C2H4, C2H6, CO2, H2O, O2, N2, CO, H2)", e.Name)
}

// UnknownFuelError reports a fuel name that is not in the pinned fuel table.
type UnknownFuelError struct {
	Name string
}

func (e *UnknownFuelError) Error() string {
	return fmt.Sprintf("unknown fuel %q (supported: CH4, C2H4, C2H6)", e.Name)
}

// TemperatureError reports a temperature outside the polynomial window.
type TemperatureError struct {
	Temperature float64
}

func (e *TemperatureError) Error() string {
	return fmt.Sprintf("temperature %.1f K is outside the polynomial window [200, 3500] K", e.Temperature)
}
