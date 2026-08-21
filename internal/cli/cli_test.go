package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCLI invokes the runner with the given arguments and returns the exit
// code plus the captured stdout and stderr.
func runCLI(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	r := &Runner{Out: &out, ErrOut: &errOut, Args: args}
	code = r.Run()
	return code, out.String(), errOut.String()
}

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "case.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestTadCommandPrintsResults(t *testing.T) {
	// The tad subcommand must print the flame temperature, the product
	// mole fractions and the atom residuals for a valid config file.
	path := writeConfig(t, `{"fuel":"CH4","equivalence_ratio":1.0,"inlet_temperature_k":298.15}`)
	code, stdout, stderr := runCLI(t, "tad", path)
	if code != 0 {
		t.Fatalf("exit %d, stderr: %s", code, stderr)
	}
	for _, want := range []string{"Adiabatic flame temperature", "Products (mole fractions)", "Atom residuals", "CO2", "N2"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q:\n%s", want, stdout)
		}
	}
	var tad float64
	if _, err := fmt.Sscanf(stdout, "Adiabatic flame temperature : %f K", &tad); err != nil {
		t.Errorf("cannot read Tad from output: %v\n%s", err, stdout)
	}
	if tad < 2000 || tad > 2400 {
		t.Errorf("printed Tad %.1f K outside [2000, 2400]", tad)
	}
}

func TestTadCommandEquivalenceRatioZero(t *testing.T) {
	// A config with phi=0 must fail loudly: exit non-zero and a message on
	// stderr, never a printed temperature.
	path := writeConfig(t, `{"fuel":"CH4","equivalence_ratio":0.0,"inlet_temperature_k":298.15}`)
	code, stdout, stderr := runCLI(t, "tad", path)
	if code == 0 {
		t.Error("phi=0 must exit non-zero")
	}
	if stderr == "" {
		t.Error("phi=0 must write an error to stderr")
	}
	if !strings.Contains(stderr, "equivalence ratio") {
		t.Errorf("stderr must explain the bad equivalence ratio, got %q", stderr)
	}
	if strings.Contains(stdout, "Adiabatic flame temperature") {
		t.Error("stdout must not contain a flame temperature for invalid input")
	}
}

func TestTadCommandInletZero(t *testing.T) {
	path := writeConfig(t, `{"fuel":"CH4","equivalence_ratio":1.0,"inlet_temperature_k":0.0}`)
	code, _, stderr := runCLI(t, "tad", path)
	if code == 0 {
		t.Error("Tin=0 must exit non-zero")
	}
	if !strings.Contains(stderr, "inlet temperature") {
		t.Errorf("stderr must explain the bad inlet temperature, got %q", stderr)
	}
}

func TestTadCommandUnknownFuel(t *testing.T) {
	path := writeConfig(t, `{"fuel":"C8H18","equivalence_ratio":1.0,"inlet_temperature_k":298.15}`)
	code, _, stderr := runCLI(t, "tad", path)
	if code == 0 {
		t.Error("unknown fuel must exit non-zero")
	}
	if !strings.Contains(stderr, "unknown fuel") {
		t.Errorf("stderr must name the unknown fuel, got %q", stderr)
	}
}

func TestTadCommandMissingFile(t *testing.T) {
	code, _, stderr := runCLI(t, "tad", "does-not-exist.json")
	if code == 0 {
		t.Error("missing file must exit non-zero")
	}
	if stderr == "" {
		t.Error("missing file must write to stderr")
	}
}

func TestTadCommandMissingArg(t *testing.T) {
	code, _, stderr := runCLI(t, "tad")
	if code != 2 {
		t.Errorf("missing argument must exit 2, got %d", code)
	}
	if stderr == "" {
		t.Error("missing argument must write to stderr")
	}
}

func TestFuelsCommand(t *testing.T) {
	code, stdout, _ := runCLI(t, "fuels")
	if code != 0 {
		t.Fatalf("fuels exit %d", code)
	}
	for _, want := range []string{"CH4", "C2H4", "C2H6"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("fuels output missing %q:\n%s", want, stdout)
		}
	}
}

func TestChecksCommandPass(t *testing.T) {
	code, stdout, stderr := runCLI(t, "checks", "CH4")
	if code != 0 {
		t.Fatalf("checks exit %d, stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "Tad(phi=1.0)") || !strings.Contains(stdout, "N2 conserved") {
		t.Errorf("checks output missing findings:\n%s", stdout)
	}
}

func TestHelpCommand(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"-h"}, {"--help"}} {
		code, stdout, _ := runCLI(t, args...)
		if code != 0 {
			t.Errorf("%v must exit 0, got %d", args, code)
		}
		if !strings.Contains(stdout, "adiab-flame tad") {
			t.Errorf("%v must show usage", args)
		}
	}
}

func TestUnknownCommand(t *testing.T) {
	code, _, stderr := runCLI(t, "explode")
	if code != 2 {
		t.Errorf("unknown command must exit 2, got %d", code)
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("stderr must reject the unknown command, got %q", stderr)
	}
}
