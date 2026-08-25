package flame

import "adiab-flame/internal/thermo"

const (
	DefaultMaxIterations   = 200
	DefaultTemperatureTolK = 0.01
	DefaultResidualTolK    = 0.5
	DefaultPressureAtm     = 1.0
	DefaultDissociation    = false
)

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

func NewConfig(fuel string, phi, tin float64) *Config {
	return (&Config{
		Fuel:             fuel,
		EquivalenceRatio: phi,
		InletTemperature: tin,
	}).WithDefaults()
}

func (c *Config) FuelName() string {
	return thermo.FormatName(c.Fuel)
}
