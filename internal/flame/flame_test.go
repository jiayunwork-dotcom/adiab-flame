package flame

import (
	"errors"
	"testing"
)

// newTestSolver returns a solver on the default registries.
func newTestSolver() *Solver {
	return NewSolver()
}

func TestMethaneStoichTadBand(t *testing.T) {
	// Pinned band: stoichiometric CH4 with a 298.15 K inlet must land in
	// the 2200 K order of magnitude. If the enthalpy forgot the formation
	// term the temperature would stay near the inlet and fail this band.
	s := newTestSolver()
	res, err := s.SolveSimple("CH4", 1.0, 298.15)
	if err != nil {
		t.Fatalf("SolveSimple(CH4, 1.0, 298.15): %v", err)
	}
	if res.Tad < 2000.0 || res.Tad > 2400.0 {
		t.Errorf("Tad = %.2f K, want in [2000, 2400] K", res.Tad)
	}
}

func TestAtomResidualsNegligible(t *testing.T) {
	// Every element must close to machine precision on the per-mole-fuel
	// basis for lean, stoichiometric and rich mixtures alike.
	s := newTestSolver()
	for _, phi := range []float64{0.8, 1.0, 1.2} {
		res, err := s.SolveSimple("CH4", phi, 298.15)
		if err != nil {
			t.Fatalf("phi=%v: %v", phi, err)
		}
		if !res.ResidualsNegligible() {
			t.Errorf("phi=%v residuals not negligible: max %.3g", phi, res.MaxResidual())
		}
	}
}

func TestEquivalenceRatioValidation(t *testing.T) {
	// A zero equivalence ratio is not a flame; the solver must reject it
	// with an error instead of dividing by zero.
	s := newTestSolver()
	for _, phi := range []float64{0.0, -0.5} {
		_, err := s.SolveSimple("CH4", phi, 298.15)
		if err == nil {
			t.Errorf("phi=%v must be rejected", phi)
		}
		var e *EquivalenceRatioError
		if !errors.As(err, &e) {
			t.Errorf("phi=%v error %T, want *EquivalenceRatioError", phi, err)
		}
	}
}

func TestInletTemperatureValidation(t *testing.T) {
	// A zero inlet temperature is outside the enthalpy model; it must be
	// rejected with a descriptive error.
	s := newTestSolver()
	for _, tin := range []float64{0.0, -100.0} {
		_, err := s.SolveSimple("CH4", 1.0, tin)
		if err == nil {
			t.Errorf("Tin=%v must be rejected", tin)
		}
		var e *InletTemperatureError
		if !errors.As(err, &e) {
			t.Errorf("Tin=%v error %T, want *InletTemperatureError", tin, err)
		}
	}
}

func TestUnknownFuelValidation(t *testing.T) {
	// A fuel outside the pinned table must be rejected before any
	// calculation runs.
	s := newTestSolver()
	for _, fuel := range []string{"O3", "C8H18", "kerosene", ""} {
		_, err := s.SolveSimple(fuel, 1.0, 298.15)
		if err == nil {
			t.Errorf("fuel %q must be rejected", fuel)
		}
	}
}

func TestLeanTadLowerThanStoich(t *testing.T) {
	// Moving lean from phi=1 must lower the flame temperature because the
	// excess air acts as a heat sink.
	s := newTestSolver()
	stoich, err := s.SolveSimple("CH4", 1.0, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	lean, err := s.SolveSimple("CH4", 0.8, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	if lean.Tad >= stoich.Tad-100.0 {
		t.Errorf("lean Tad %.2f K must be at least 100 K below stoich Tad %.2f K", lean.Tad, stoich.Tad)
	}
}

func TestRichTadLowerThanStoich(t *testing.T) {
	// Moving rich from phi=1 must lower the flame temperature; the
	// unburned fuel dilutes the products.
	s := newTestSolver()
	stoich, err := s.SolveSimple("CH4", 1.0, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	rich, err := s.SolveSimple("CH4", 1.2, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	if rich.Tad >= stoich.Tad-30.0 {
		t.Errorf("rich Tad %.2f K must be below stoich Tad %.2f K", rich.Tad, stoich.Tad)
	}
}

func TestHotInletTadRisesSlowly(t *testing.T) {
	// Preheating the inlet must raise Tad, but by less than the inlet
	// delta because the hot products absorb more sensible heat.
	s := newTestSolver()
	cold, err := s.SolveSimple("CH4", 1.0, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	hot, err := s.SolveSimple("CH4", 1.0, 600.0)
	if err != nil {
		t.Fatal(err)
	}
	deltaIn := 600.0 - 298.15
	deltaOut := hot.Tad - cold.Tad
	if deltaOut <= 0 {
		t.Errorf("hot inlet must raise Tad, got delta %.2f K", deltaOut)
	}
	if deltaOut >= deltaIn {
		t.Errorf("Tad delta %.2f K must stay below inlet delta %.2f K", deltaOut, deltaIn)
	}
}

func TestNitrogenConserved(t *testing.T) {
	// N2 is inert in this model: the product stream must carry exactly the
	// nitrogen that arrived with the air.
	s := newTestSolver()
	for _, phi := range []float64{0.8, 1.0, 1.2} {
		res, err := s.SolveSimple("CH4", phi, 298.15)
		if err != nil {
			t.Fatalf("phi=%v: %v", phi, err)
		}
		inlet, outlet := res.NitrogenBalance()
		if diff := inlet - outlet; diff > 1e-9 || diff < -1e-9 {
			t.Errorf("phi=%v N2 in %.6g out %.6g, want equal", phi, inlet, outlet)
		}
	}
}
