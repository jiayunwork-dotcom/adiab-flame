package cli

import (
	"fmt"
	"strings"

	"adiab-flame/internal/flame"
)

func FormatResult(res *flame.Result) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Adiabatic flame temperature : %8.2f K\n", res.Tad)
	fmt.Fprintf(&b, "Fuel                        : %s\n", res.Fuel)
	fmt.Fprintf(&b, "Equivalence ratio           : %6.3f\n", res.EquivalenceRatio)
	fmt.Fprintf(&b, "Inlet temperature           : %8.2f K\n", res.InletTemperature)
	fmt.Fprintf(&b, "Pressure                    : %6.3f atm\n", res.PressureAtm)
	fmt.Fprintf(&b, "Iterations                  : %d\n", res.Iterations)
	fmt.Fprintf(&b, "Reactant enthalpy           : %10.3f kJ/mol fuel\n", res.ReactantEnthalpy)
	fmt.Fprintf(&b, "Product enthalpy            : %10.3f kJ/mol fuel\n", res.ProductEnthalpy)
	fmt.Fprintf(&b, "Dissociation refinement     : %v\n", res.DissociationUsed())
	b.WriteString("\n")

	b.WriteString("Products (mole fractions):\n")
	names := displayNames(res)
	for _, name := range names {
		x := res.MoleFraction(name)
		if x == 0 && !res.Products.SpeciesPresent(name) {
			continue
		}
		fmt.Fprintf(&b, "  %-5s %10.5f\n", name, x)
	}
	b.WriteString("\n")

	b.WriteString("Atom residuals (mol per mol fuel):\n")
	for _, r := range res.Residuals {
		fmt.Fprintf(&b, "  %-2s  in %10.6f  out %10.6f  delta %+.3e\n", r.Element, r.Inlet, r.Outlet, r.Delta)
	}
	if res.ResidualsNegligible() {
		fmt.Fprintf(&b, "  all residuals below %.0e\n", flame.ResidualThreshold)
	} else {
		fmt.Fprintf(&b, "  residual threshold %.0e exceeded\n", flame.ResidualThreshold)
	}
	return b.String()
}

func displayNames(res *flame.Result) []string {
	out := make([]string, 0, 8)
	for _, name := range []string{"CO2", "H2O", "O2", "N2", "CO", "H2"} {
		out = append(out, name)
	}
	if res.Products.Fuel != "" {
		found := false
		for _, n := range out {
			if n == res.Products.Fuel {
				found = true
				break
			}
		}
		if !found {
			out = append(out, res.Products.Fuel)
		}
	}
	return out
}

func FormatCompact(res *flame.Result) string {
	return fmt.Sprintf("Tad=%.2f K  phi=%.3f  fuel=%s  maxres=%.3g",
		res.Tad, res.EquivalenceRatio, res.Fuel, res.MaxResidual())
}
