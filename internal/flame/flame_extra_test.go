package flame

import (
	"strings"
	"testing"
)

func TestRichUnburnedFuel(t *testing.T) {
	s := newTestSolver()
	res, err := s.SolveSimple("CH4", 1.2, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	fuelAmount := res.ProductAmount("CH4")
	want := 1.0 - 1.0/1.2
	if diff := fuelAmount - want; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("unburned CH4 = %.6g mol, want %.6g", fuelAmount, want)
	}
	if res.Products.Fuel != "CH4" {
		t.Errorf("product stream fuel tag = %q, want CH4", res.Products.Fuel)
	}
	if res.ProductAmount("O2") != 0 {
		t.Errorf("rich products must not contain leftover O2, got %.3g", res.ProductAmount("O2"))
	}
}

func TestDissociationFormsCO(t *testing.T) {
	s := newTestSolver()
	cfg := NewConfig("CH4", 1.0, 298.15)
	cfg.Dissociation = true
	res, err := s.Solve(cfg)
	if err != nil {
		t.Fatal(err)
	}
	co := res.MoleFraction("CO")
	h2 := res.MoleFraction("H2")
	if co <= 0.001 {
		t.Errorf("CO mole fraction %.5f, want > 0.001 with dissociation", co)
	}
	if h2 <= 0.0001 {
		t.Errorf("H2 mole fraction %.5f, want > 0.0001 with dissociation", h2)
	}
	if co > 0.1 || h2 > 0.1 {
		t.Errorf("CO/H2 fractions too large: CO=%.4f H2=%.4f", co, h2)
	}
	if !res.DissociationUsed() {
		t.Error("result must report dissociation was used")
	}
}

func TestDissociationLowersTad(t *testing.T) {
	s := newTestSolver()
	plain, err := s.SolveSimple("CH4", 1.0, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	cfg := NewConfig("CH4", 1.0, 298.15)
	cfg.Dissociation = true
	dissoc, err := s.Solve(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if dissoc.Tad >= plain.Tad {
		t.Errorf("dissociation Tad %.2f K must be below plain Tad %.2f K", dissoc.Tad, plain.Tad)
	}
}

func TestEthaneStoichTad(t *testing.T) {
	s := newTestSolver()
	res, err := s.SolveSimple("C2H6", 1.0, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	if res.Tad < 2100.0 || res.Tad > 2600.0 {
		t.Errorf("C2H6 stoich Tad = %.2f K, want in [2100, 2600] K", res.Tad)
	}
	if !res.ResidualsNegligible() {
		t.Errorf("C2H6 residuals not negligible: %.3g", res.MaxResidual())
	}
}

func TestEthyleneStoichTad(t *testing.T) {
	s := newTestSolver()
	res, err := s.SolveSimple("C2H4", 1.0, 298.15)
	if err != nil {
		t.Fatal(err)
	}
	if res.Tad < 2300.0 || res.Tad > 2750.0 {
		t.Errorf("C2H4 stoich Tad = %.2f K, want in [2300, 2750] K", res.Tad)
	}
	if !res.ResidualsNegligible() {
		t.Errorf("C2H4 residuals not negligible: %.3g", res.MaxResidual())
	}
}

func TestCrossChecksAllPass(t *testing.T) {
	s := newTestSolver()
	for _, fuel := range []string{"CH4", "C2H4", "C2H6"} {
		rep, err := RunCrossChecks(s, fuel, 298.15)
		if err != nil {
			t.Fatalf("%s: %v", fuel, err)
		}
		if !rep.AllPass() {
			t.Errorf("%s cross rules failed: %s", fuel, strings.Join(rep.Failures, "; "))
		}
	}
}

func TestConfigFromJSONValidation(t *testing.T) {
	bad := `{"fuel": "CH4", "equivalence_ratio": 1.0, "inlet_temperature_k": 298.15, "tempreature": 500}`
	_, err := ConfigFromJSON([]byte(bad))
	if err == nil {
		t.Error("unknown JSON field must be rejected")
	} else if !strings.Contains(err.Error(), "tempreature") {
		t.Errorf("error must name the unknown field, got %v", err)
	}
	ok := `{"fuel": "CH4", "equivalence_ratio": 1.0, "inlet_temperature_k": 298.15}`
	cfg, err := ConfigFromJSON([]byte(ok))
	if err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	if cfg.EquivalenceRatio != 1.0 || cfg.InletTemperature != 298.15 {
		t.Errorf("config parsed wrongly: %+v", cfg)
	}
}
