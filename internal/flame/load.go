package flame

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type jsonConfig struct {
	Fuel             string   `json:"fuel"`
	EquivalenceRatio *float64 `json:"equivalence_ratio"`
	InletTemperature *float64 `json:"inlet_temperature_k"`
	PressureAtm      *float64 `json:"pressure_atm,omitempty"`
	Dissociation     *bool    `json:"dissociation,omitempty"`
	MaxIterations    *int     `json:"max_iterations,omitempty"`
}

var allowedJSONFields = map[string]bool{
	"fuel":                true,
	"equivalence_ratio":   true,
	"inlet_temperature_k": true,
	"pressure_atm":        true,
	"dissociation":        true,
	"max_iterations":      true,
}

func ConfigFromJSON(data []byte) (*Config, error) {
	if err := rejectUnknownFields(data); err != nil {
		return nil, err
	}
	var jc jsonConfig
	if err := json.Unmarshal(data, &jc); err != nil {
		return nil, fmt.Errorf("parse config JSON: %w", err)
	}
	return buildConfig(&jc)
}

func ConfigFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return ConfigFromJSON(data)
}

func buildConfig(jc *jsonConfig) (*Config, error) {
	if jc.Fuel == "" {
		return nil, &ConfigError{Field: "fuel", Message: "missing fuel name"}
	}
	cfg := &Config{Fuel: jc.Fuel}
	cfg.WithDefaults()

	if jc.EquivalenceRatio == nil {
		return nil, &ConfigError{Field: "equivalence_ratio", Message: "missing equivalence ratio"}
	}
	cfg.EquivalenceRatio = *jc.EquivalenceRatio

	if jc.InletTemperature == nil {
		return nil, &ConfigError{Field: "inlet_temperature_k", Message: "missing inlet temperature"}
	}
	cfg.InletTemperature = *jc.InletTemperature

	if jc.PressureAtm != nil {
		cfg.PressureAtm = *jc.PressureAtm
	}
	if jc.Dissociation != nil {
		cfg.Dissociation = *jc.Dissociation
	}
	if jc.MaxIterations != nil {
		cfg.MaxIterations = *jc.MaxIterations
	}
	return cfg, nil
}

func rejectUnknownFields(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse config JSON: %w", err)
	}
	for key := range raw {
		if !allowedJSONFields[key] {
			return &ConfigError{Field: key, Message: "unknown JSON field"}
		}
	}
	return nil
}
