package cli

import (
	"fmt"
	"io"
	"os"
)

// Runner executes the command line interface. Output goes to out and
// errors to errOut; both default to stdout and stderr when nil.
type Runner struct {
	Out    io.Writer
	ErrOut io.Writer
	Args   []string
}

// Run dispatches the subcommand and returns an exit code. Errors that
// belong to the user (bad input, non-convergence) are printed on stderr
// and mapped to a non-zero exit; usage mistakes exit with code 2.
func (r *Runner) Run() int {
	if r.Out == nil {
		r.Out = os.Stdout
	}
	if r.ErrOut == nil {
		r.ErrOut = os.Stderr
	}
	if len(r.Args) == 0 {
		fmt.Fprint(r.ErrOut, usage)
		return 2
	}

	switch r.Args[0] {
	case "tad":
		return r.runTad(r.Args[1:])
	case "checks":
		return r.runChecks(r.Args[1:])
	case "fuels":
		return r.runFuels()
	case "help", "-h", "--help":
		fmt.Fprint(r.Out, usage)
		return 0
	default:
		fmt.Fprintf(r.ErrOut, "adiab-flame: unknown command %q\n\n%s", r.Args[0], usage)
		return 2
	}
}

// RunCommand is the package level entry point used by main.
func RunCommand(args []string) int {
	return (&Runner{Args: args}).Run()
}

// runTad implements the `tad` subcommand: solve one config file and print
// the result.
func (r *Runner) runTad(args []string) int {
	path, dissoc, err := parseTadArgs(args)
	if err != nil {
		fmt.Fprintf(r.ErrOut, "adiab-flame tad: %v\n", err)
		return 2
	}
	res, err := SolveFile(path, dissoc)
	if err != nil {
		fmt.Fprintf(r.ErrOut, "adiab-flame tad: %v\n", err)
		return 1
	}
	fmt.Fprint(r.Out, FormatResult(res))
	return 0
}

// runChecks implements the `checks` subcommand: verify the pinned cross
// rules for one fuel.
func (r *Runner) runChecks(args []string) int {
	fuel, inlet, err := parseChecksArgs(args)
	if err != nil {
		fmt.Fprintf(r.ErrOut, "adiab-flame checks: %v\n", err)
		return 2
	}
	rep, err := RunCrossChecks(fuel, inlet)
	if err != nil {
		fmt.Fprintf(r.ErrOut, "adiab-flame checks: %v\n", err)
		return 1
	}
	fmt.Fprint(r.Out, rep.Summary())
	if !rep.AllPass() {
		fmt.Fprintf(r.ErrOut, "adiab-flame checks: %d cross-rule failure(s)\n", len(rep.Failures))
		return 1
	}
	return 0
}

// runFuels implements the `fuels` subcommand: list supported fuels.
func (r *Runner) runFuels() int {
	fmt.Fprint(r.Out, FuelList())
	return 0
}
