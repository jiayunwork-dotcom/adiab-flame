package flame

import "fmt"

// ConfigError reports a missing or malformed configuration field.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config: field %q: %s", e.Field, e.Message)
}

// UnknownFuelInputError reports a fuel name the solver does not know.
type UnknownFuelInputError struct {
	Fuel  string
	Known []string
}

func (e *UnknownFuelInputError) Error() string {
	return fmt.Sprintf("unknown fuel %q (supported fuels: %v)", e.Fuel, e.Known)
}

// EquivalenceRatioError reports a non-positive equivalence ratio. The
// equivalence ratio must be strictly positive; phi = 0 means no fuel and
// negative ratios are not physical.
type EquivalenceRatioError struct {
	Ratio float64
}

func (e *EquivalenceRatioError) Error() string {
	return fmt.Sprintf("equivalence ratio must be > 0, got %v", e.Ratio)
}

// InletTemperatureError reports a non-positive inlet temperature. The
// enthalpy polynomials are not defined at or below zero kelvin.
type InletTemperatureError struct {
	Temperature float64
}

func (e *InletTemperatureError) Error() string {
	return fmt.Sprintf("inlet temperature must be > 0 K, got %v", e.Temperature)
}

// NonConvergenceError reports a temperature iteration that used up its
// iteration budget without closing the enthalpy balance.
type NonConvergenceError struct {
	MaxIterations    int
	LastBracketWidth float64
}

func (e *NonConvergenceError) Error() string {
	return fmt.Sprintf("temperature iteration did not converge after %d iterations (bracket width %.3f K)", e.MaxIterations, e.LastBracketWidth)
}

// BracketError reports that the enthalpy balance did not straddle zero on
// the search bracket, so bisection cannot proceed.
type BracketError struct {
	Low     float64
	High    float64
	Message string
}

func (e *BracketError) Error() string {
	return fmt.Sprintf("temperature bracket [%.2f, %.2f]: %s", e.Low, e.High, e.Message)
}

// DissociationError reports a failure inside the optional equilibrium
// refinement.
type DissociationError struct {
	Temperature float64
	Reason      string
}

func (e *DissociationError) Error() string {
	return fmt.Sprintf("dissociation solve failed at %.1f K: %s", e.Temperature, e.Reason)
}
