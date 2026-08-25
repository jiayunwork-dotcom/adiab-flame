package flame

import "fmt"

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config: field %q: %s", e.Field, e.Message)
}

type UnknownFuelInputError struct {
	Fuel  string
	Known []string
}

func (e *UnknownFuelInputError) Error() string {
	return fmt.Sprintf("unknown fuel %q (supported fuels: %v)", e.Fuel, e.Known)
}

type EquivalenceRatioError struct {
	Ratio float64
}

func (e *EquivalenceRatioError) Error() string {
	return fmt.Sprintf("equivalence ratio must be > 0, got %v", e.Ratio)
}

type InletTemperatureError struct {
	Temperature float64
}

func (e *InletTemperatureError) Error() string {
	return fmt.Sprintf("inlet temperature must be > 0 K, got %v", e.Temperature)
}

type NonConvergenceError struct {
	MaxIterations    int
	LastBracketWidth float64
}

func (e *NonConvergenceError) Error() string {
	return fmt.Sprintf("temperature iteration did not converge after %d iterations (bracket width %.3f K)", e.MaxIterations, e.LastBracketWidth)
}

type BracketError struct {
	Low     float64
	High    float64
	Message string
}

func (e *BracketError) Error() string {
	return fmt.Sprintf("temperature bracket [%.2f, %.2f]: %s", e.Low, e.High, e.Message)
}

type DissociationError struct {
	Temperature float64
	Reason      string
}

func (e *DissociationError) Error() string {
	return fmt.Sprintf("dissociation solve failed at %.1f K: %s", e.Temperature, e.Reason)
}
