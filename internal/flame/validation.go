package flame

import (
	"fmt"
	"math"

	"adiab-flame/internal/thermo"
)

func (c *Config) Validate(registry *thermo.Registry, fuels *thermo.FuelRegistry) error {
	cfg := c.WithDefaults()

	if cfg.Fuel == "" {
		return &ConfigError{Field: "fuel", Message: "missing fuel name"}
	}
	if !fuels.Has(cfg.Fuel) {
		return &UnknownFuelInputError{Fuel: cfg.Fuel, Known: fuels.Names()}
	}

	if math.IsNaN(cfg.EquivalenceRatio) || math.IsInf(cfg.EquivalenceRatio, 0) {
		return &ConfigError{Field: "equivalence_ratio", Message: "must be a finite number"}
	}
	if cfg.EquivalenceRatio <= 0 {
		return &EquivalenceRatioError{Ratio: cfg.EquivalenceRatio}
	}

	if math.IsNaN(cfg.InletTemperature) || math.IsInf(cfg.InletTemperature, 0) {
		return &ConfigError{Field: "inlet_temperature_k", Message: "must be a finite number"}
	}
	if cfg.InletTemperature <= 0 {
		return &InletTemperatureError{Temperature: cfg.InletTemperature}
	}
	if !thermo.IsValidTemperature(cfg.InletTemperature) {
		return &ConfigError{
			Field:   "inlet_temperature_k",
			Message: fmt.Sprintf("%.1f K is outside the polynomial window [200, 3500] K", cfg.InletTemperature),
		}
	}

	if cfg.PressureAtm <= 0 {
		return &ConfigError{Field: "pressure_atm", Message: "must be positive"}
	}
	if cfg.MaxIterations <= 0 {
		return &ConfigError{Field: "max_iterations", Message: "must be positive"}
	}
	return nil
}

func (c *Config) MustResolveFuel(fuels *thermo.FuelRegistry) *thermo.Fuel {
	f, err := fuels.Lookup(c.Fuel)
	if err != nil {
		panic(err)
	}
	return f
}
