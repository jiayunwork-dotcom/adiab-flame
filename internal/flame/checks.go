package flame

import (
	"fmt"
	"math"
)

// CrossCheckReport summarises the pinned cross rules for one fuel and one
// inlet temperature: the flame temperatures at three equivalence ratios,
// the response to a hotter inlet, and the conservation checks at the
// stoichiometric point.
type CrossCheckReport struct {
	Fuel         string
	InletK       float64
	HotInletK    float64
	TadStoichK   float64
	TadLeanK     float64
	TadRichK     float64
	TadHotInletK float64
	MaxResidual  float64
	N2Inlet      float64
	N2Outlet     float64
	Checks       []string // human readable findings, pass and fail
	Failures     []string // findings that violate a cross rule
}

// RunCrossChecks evaluates every cross rule for a fuel and inlet
// temperature. It is the source of truth the `checks` CLI command prints.
func RunCrossChecks(s *Solver, fuel string, inletK float64) (*CrossCheckReport, error) {
	report := &CrossCheckReport{Fuel: fuel, InletK: inletK, HotInletK: inletK + 300.0}

	stoich, err := s.SolveSimple(fuel, 1.0, inletK)
	if err != nil {
		return nil, err
	}
	report.TadStoichK = stoich.Tad
	report.MaxResidual = stoich.MaxResidual()

	lean, err := s.SolveSimple(fuel, 0.8, inletK)
	if err != nil {
		return nil, err
	}
	report.TadLeanK = lean.Tad

	rich, err := s.SolveSimple(fuel, 1.2, inletK)
	if err != nil {
		return nil, err
	}
	report.TadRichK = rich.Tad

	hot, err := s.SolveSimple(fuel, 1.0, report.HotInletK)
	if err != nil {
		return nil, err
	}
	report.TadHotInletK = hot.Tad

	report.N2Inlet, report.N2Outlet = stoich.NitrogenBalance()
	report.evaluate()
	return report, nil
}

// evaluate scores every cross rule and fills the Checks/Failures lists.
func (r *CrossCheckReport) evaluate() {
	r.Checks = nil
	r.Failures = nil

	r.add("stoichiometric Tad %.1f K", r.TadStoichK)
	r.add("lean (phi=0.8) Tad %.1f K", r.TadLeanK)
	r.add("rich (phi=1.2) Tad %.1f K", r.TadRichK)

	if r.TadStoichK > r.TadLeanK+50 {
		r.add("phi=1 Tad exceeds phi=0.8 Tad by %.1f K", r.TadStoichK-r.TadLeanK)
	} else {
		r.fail("phi=1 Tad (%.1f K) must exceed phi=0.8 Tad (%.1f K)", r.TadStoichK, r.TadLeanK)
	}

	if r.TadStoichK > r.TadRichK+30 {
		r.add("phi=1 Tad exceeds phi=1.2 Tad by %.1f K", r.TadStoichK-r.TadRichK)
	} else {
		r.fail("phi=1 Tad (%.1f K) must exceed phi=1.2 Tad (%.1f K)", r.TadStoichK, r.TadRichK)
	}

	deltaIn := r.HotInletK - r.InletK
	deltaOut := r.TadHotInletK - r.TadStoichK
	if r.TadHotInletK > r.TadStoichK {
		r.add("hotter inlet raises Tad by %.1f K (inlet rose %.1f K)", deltaOut, deltaIn)
	} else {
		r.fail("raising the inlet temperature must raise Tad")
	}
	if deltaOut < deltaIn {
		r.add("Tad rise %.1f K stays below the inlet rise %.1f K", deltaOut, deltaIn)
	} else {
		r.fail("Tad rise %.1f K must stay below the inlet rise %.1f K", deltaOut, deltaIn)
	}

	if r.MaxResidual < ResidualThreshold {
		r.add("atom residuals max %.3g below threshold %.3g", r.MaxResidual, ResidualThreshold)
	} else {
		r.fail("atom residuals max %.3g must stay below threshold %.3g", r.MaxResidual, ResidualThreshold)
	}

	if math.Abs(r.N2Inlet-r.N2Outlet) < 1e-9 {
		r.add("N2 conserved: %.6g mol in, %.6g mol out", r.N2Inlet, r.N2Outlet)
	} else {
		r.fail("N2 not conserved: %.6g mol in, %.6g mol out", r.N2Inlet, r.N2Outlet)
	}
}

func (r *CrossCheckReport) add(format string, args ...any) {
	r.Checks = append(r.Checks, fmt.Sprintf("OK   "+format, args...))
}

func (r *CrossCheckReport) fail(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	r.Checks = append(r.Checks, "FAIL "+msg)
	r.Failures = append(r.Failures, msg)
}

// AllPass reports whether every cross rule held.
func (r *CrossCheckReport) AllPass() bool {
	return len(r.Failures) == 0
}

// Summary renders the report for the CLI in a compact block.
func (r *CrossCheckReport) Summary() string {
	out := ""
	out += fmt.Sprintf("fuel %s, inlet %.2f K\n", r.Fuel, r.InletK)
	out += fmt.Sprintf("  Tad(phi=1.0) = %.2f K\n", r.TadStoichK)
	out += fmt.Sprintf("  Tad(phi=0.8) = %.2f K\n", r.TadLeanK)
	out += fmt.Sprintf("  Tad(phi=1.2) = %.2f K\n", r.TadRichK)
	out += fmt.Sprintf("  Tad(inlet %.0f K) = %.2f K\n", r.HotInletK, r.TadHotInletK)
	out += fmt.Sprintf("  max atom residual = %.3g\n", r.MaxResidual)
	out += fmt.Sprintf("  N2 in = %.6g mol, N2 out = %.6g mol\n", r.N2Inlet, r.N2Outlet)
	for _, c := range r.Checks {
		out += "  " + c + "\n"
	}
	return out
}
