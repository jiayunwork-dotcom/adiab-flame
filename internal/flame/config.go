package flame

import "adiab-flame/internal/thermo"

// Default values pinned for every run unless overridden by JSON.
const (
	// DefaultMaxIterations caps the temperature bisection.
	DefaultMaxIterations = 200
	// DefaultTemperatureTolK stops the bisection once the bracket is this narrow.
	DefaultTemperatureTolK = 0.01
	// DefaultResidualTolK bounds the enthalpy residual converted to kelvin.
	DefaultResidualTolK = 0.5
	// DefaultPressureAtm is the constant pressure of the flame.
	DefaultPressureAtm = 1.0
	// DefaultDissociation switches the optional CO/H2 equilibrium off.
	DefaultDissociation = false
)

// Config describes one adiabatic flame calculation: fuel identity, the
// equivalence ratio, the inlet temperature and the modelling switches.
type Config struct {
	Fuel             string
	EquivalenceRatio float64
	InletTemperature float64
	PressureAtm      float64
	Dissociation     bool
	MaxIterations    int
	TemperatureTolK  float64
	ResidualTolK     float64
}

// WithDefaults fills any zero-valued tuning field with the pinned default
// so callers can construct a Config with only the physical inputs set.
func (c *Config) WithDefaults() *Config {
	out := *c
	if out.MaxIterations <= 0 {
		out.MaxIterations = DefaultMaxIterations
	}
	if out.TemperatureTolK <= 0 {
		out.TemperatureTolK = DefaultTemperatureTolK
	}
	if out.ResidualTolK <= 0 {
		out.ResidualTolK = DefaultResidualTolK
	}
	if out.PressureAtm <= 0 {
		out.PressureAtm = DefaultPressureAtm
	}
	return &out
}

// NewConfig returns a Config with the physical inputs set and every tuning
// field on its default.
func NewConfig(fuel string, phi, tin float64) *Config {
	return (&Config{
		Fuel:             fuel,
		EquivalenceRatio: phi,
		InletTemperature: tin,
	}).WithDefaults()
}

// FuelName returns the resolved fuel name used by the solver.
func (c *Config) FuelName() string {
	return thermo.FormatName(c.Fuel)
}
