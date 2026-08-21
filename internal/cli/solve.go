package cli

import (
	"errors"
	"fmt"
	"strings"

	"adiab-flame/internal/flame"
	"adiab-flame/internal/thermo"
)

// SolveFile loads a config JSON file and runs the solver with the optional
// dissociation switch applied.
func SolveFile(path string, dissoc bool) (*flame.Result, error) {
	cfg, err := flame.ConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	if dissoc {
		cfg.Dissociation = true
	}
	return flame.NewSolver().Solve(cfg)
}

// RunCrossChecks executes the pinned cross rules for a fuel at an inlet
// temperature.
func RunCrossChecks(fuel string, inlet float64) (*flame.CrossCheckReport, error) {
	if inlet <= 0 {
		return nil, &InputError{Message: fmt.Sprintf("inlet temperature must be > 0 K, got %v", inlet)}
	}
	return flame.RunCrossChecks(flame.NewSolver(), fuel, inlet)
}

// FuelList renders the supported fuels with their formation enthalpies.
func FuelList() string {
	reg := thermo.DefaultFuelRegistry
	lines := make([]string, 0, reg.Count()+1)
	lines = append(lines, "supported fuels (enthalpy of formation at 298.15 K):")
	for _, name := range reg.Names() {
		spec, err := thermo.NewSpecies(name)
		if err != nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("  %-6s  dhf = %8.3f kJ/mol", name, spec.Formation))
	}
	return strings.Join(lines, "\n") + "\n"
}

// parseTadArgs handles the `tad <file> [--dissoc]` argument shape. It also
// tolerates the flag before the file name.
func parseTadArgs(args []string) (path string, dissoc bool, err error) {
	dissoc = false
	for _, a := range args {
		switch a {
		case "--dissoc", "-d":
			dissoc = true
		case "":
			continue
		default:
			if strings.HasPrefix(a, "-") {
				return "", false, &InputError{Message: fmt.Sprintf("unknown flag %q", a)}
			}
			if path != "" {
				return "", false, &InputError{Message: "expected exactly one config file"}
			}
			path = a
		}
	}
	if path == "" {
		return "", false, &InputError{Message: "missing config file path"}
	}
	return path, dissoc, nil
}

// parseChecksArgs handles the `checks <fuel> [--inlet K]` argument shape.
func parseChecksArgs(args []string) (fuel string, inlet float64, err error) {
	inlet = 298.15
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--inlet" || a == "-t":
			if i+1 >= len(args) {
				return "", 0, &InputError{Message: "--inlet needs a value"}
			}
			i++
			v, perr := parseKelvin(args[i])
			if perr != nil {
				return "", 0, perr
			}
			inlet = v
		case strings.HasPrefix(a, "--inlet="):
			v, perr := parseKelvin(strings.TrimPrefix(a, "--inlet="))
			if perr != nil {
				return "", 0, perr
			}
			inlet = v
		case strings.HasPrefix(a, "-") && a != "-":
			return "", 0, &InputError{Message: fmt.Sprintf("unknown flag %q", a)}
		default:
			if fuel != "" {
				return "", 0, &InputError{Message: "expected exactly one fuel name"}
			}
			fuel = a
		}
	}
	if fuel == "" {
		return "", 0, &InputError{Message: "missing fuel name"}
	}
	return fuel, inlet, nil
}

func parseKelvin(s string) (float64, error) {
	var v float64
	if _, e := fmt.Sscanf(s, "%g", &v); e != nil {
		return 0, &InputError{Message: fmt.Sprintf("%q is not a temperature in kelvin", s)}
	}
	return v, nil
}

// InputError is a user-facing command line error.
type InputError struct {
	Message string
}

func (e *InputError) Error() string {
	return e.Message
}

// IsInputError matches every error produced by argument parsing.
func IsInputError(err error) bool {
	var ie *InputError
	return errors.As(err, &ie)
}
