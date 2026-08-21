package thermo

import (
	"math"
	"testing"
)

// formationEnthalpies holds the pinned reference values the tests compare
// against. They are independent of the implementation so the assertion
// catches both coefficient-table drift and a missing formation term.
var formationEnthalpies = map[string]float64{
	"CO2":  -393.522,
	"H2O":  -241.826,
	"O2":   0.0,
	"N2":   0.0,
	"CH4":  -74.873,
	"C2H4": 52.467,
	"C2H6": -84.684,
	"CO":   -110.527,
	"H2":   0.0,
}

func TestEnthalpyFormationIncluded(t *testing.T) {
	// At the reference temperature the enthalpy must equal the enthalpy of
	// formation; a cp-only enthalpy would evaluate to zero here.
	for name, want := range formationEnthalpies {
		spec, err := NewSpecies(name)
		if err != nil {
			t.Fatalf("species %q: %v", name, err)
		}
		got := spec.Enthalpy(ReferenceTemperature)
		if math.Abs(got-want) > 1e-3 {
			t.Errorf("Enthalpy(%s, %v) = %.3f kJ/mol, want %.3f (formation term included)", name, ReferenceTemperature, got, want)
		}
	}
}

func TestEnthalpySensibleMatchesNASAIntegration(t *testing.T) {
	// The thermal part h - dhf must equal the analytic integral of cp over
	// the same polynomial range, verified at a point in each range.
	spec, _ := NewSpecies("CO2")
	for _, tK := range []float64{300.0, 600.0, 1200.0, 2200.0} {
		shift := spec.Enthalpy(tK) - spec.Formation
		sens := SensibleEnthalpy(spec, tK) / 1000.0
		if math.Abs(shift-sens) > 1e-9 {
			t.Errorf("CO2 at %.0f K: enthalpy shift %.6f != sensible %.6f kJ/mol", tK, shift, sens)
		}
	}
}

func TestCappacityBaseline(t *testing.T) {
	// Reference heat capacities at 298.15 K in J/(mol K); values are from
	// standard tables, not from the polynomial under test.
	ref := map[string]float64{
		"N2": 29.12, "O2": 29.38, "CO2": 37.13, "H2O": 33.59,
		"CH4": 35.69, "C2H4": 42.87, "C2H6": 52.48, "CO": 29.14, "H2": 27.97,
	}
	for name, want := range ref {
		spec, err := NewSpecies(name)
		if err != nil {
			t.Fatalf("species %q: %v", name, err)
		}
		got := spec.Cappacity(ReferenceTemperature)
		if math.Abs(got-want) > 0.5 {
			t.Errorf("Cappacity(%s, 298.15) = %.2f J/(mol K), want ~%.2f", name, got, want)
		}
	}
}

func TestEnthalpyMonotonicInTemperature(t *testing.T) {
	// Enthalpy must rise with temperature because cp is positive, which is
	// what makes the flame temperature root unique.
	spec, _ := NewSpecies("N2")
	prev := spec.Enthalpy(300.0)
	for _, tK := range []float64{500.0, 800.0, 1200.0, 1800.0, 2500.0, 3200.0} {
		cur := spec.Enthalpy(tK)
		if cur <= prev {
			t.Errorf("N2 enthalpy not monotonic: h(%.0f)=%.4f <= h(prev)=%.4f", tK, cur, prev)
		}
		prev = cur
	}
}

func TestKpDissociationOrder(t *testing.T) {
	// Both dissociation constants must rise with temperature (endothermic)
	// and CO2 dissociation must sit above H2O dissociation at every point.
	prev1, prev2 := 0.0, 0.0
	for _, tK := range []float64{1500.0, 2000.0, 2500.0, 3000.0} {
		k1, err := KpCO2Dissociation(tK)
		if err != nil {
			t.Fatalf("Kp1 at %.0f K: %v", tK, err)
		}
		k2, err := KpH2ODissociation(tK)
		if err != nil {
			t.Fatalf("Kp2 at %.0f K: %v", tK, err)
		}
		if k1 <= 0 || k2 <= 0 || math.IsInf(k1, 0) || math.IsInf(k2, 0) {
			t.Errorf("Kp must be positive and finite at %.0f K, got %.3e / %.3e", tK, k1, k2)
		}
		if k1 <= prev1 || k2 <= prev2 {
			t.Errorf("Kp must rise with temperature at %.0f K", tK)
		}
		if k1 <= k2 {
			t.Errorf("CO2 dissociation (%.3e) must exceed H2O dissociation (%.3e) at %.0f K", k1, k2, tK)
		}
		prev1, prev2 = k1, k2
	}
}

func TestKpKnownOrderOfMagnitude(t *testing.T) {
	// Literature order of magnitude: around 1e-3 at 2000 K for CO2
	// dissociation and a few times 1e-4 for H2O dissociation.
	k1, _ := KpCO2Dissociation(2000.0)
	k2, _ := KpH2ODissociation(2000.0)
	if k1 < 5e-4 || k1 > 5e-3 {
		t.Errorf("Kp(CO2, 2000 K) = %.3e, want in [5e-4, 5e-3]", k1)
	}
	if k2 < 1e-4 || k2 > 1e-3 {
		t.Errorf("Kp(H2O, 2000 K) = %.3e, want in [1e-4, 1e-3]", k2)
	}
}

func TestAirMolarComposition(t *testing.T) {
	if AirOxygenFraction+AirNitrogenFraction != 1.0 {
		t.Errorf("air fractions must sum to 1, got %v + %v", AirOxygenFraction, AirNitrogenFraction)
	}
	if math.Abs(AirN2PerO2-79.0/21.0) > 1e-12 {
		t.Errorf("N2/O2 ratio = %v, want 79/21", AirN2PerO2)
	}
}

func TestRegistryLookup(t *testing.T) {
	reg := NewRegistry()
	if !reg.Has("co2") {
		t.Error("registry must accept lowercase species names")
	}
	if !reg.Has("CH4") {
		t.Error("registry must resolve CH4")
	}
	if reg.Has("O3") {
		t.Error("registry must not invent unknown species")
	}
	if _, err := reg.Lookup("O3"); err == nil {
		t.Error("Lookup(O3) must fail")
	}
	if reg.Count() != len(speciesTable) {
		t.Errorf("registry count %d != table count %d", reg.Count(), len(speciesTable))
	}
}

func TestFuelRegistry(t *testing.T) {
	reg := NewFuelRegistry()
	if !reg.Has("ch4") {
		t.Error("fuel registry must resolve CH4 case-insensitively")
	}
	if reg.Has("O2") {
		t.Error("O2 is not a fuel")
	}
	fuel, err := reg.Lookup("C2H4")
	if err != nil {
		t.Fatalf("Lookup(C2H4): %v", err)
	}
	if fuel.Carbon != 2 || fuel.Hydrogen != 4 {
		t.Errorf("C2H4 atoms = C%d H%d, want C2 H4", int(fuel.Carbon), int(fuel.Hydrogen))
	}
}
