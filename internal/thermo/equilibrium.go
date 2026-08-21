package thermo

import "math"

// EquilibriumReaction describes one gas phase dissociation reaction written
// in the form A -> B + 0.5*C, with a stoichiometric half coefficient that
// the equilibrium expression must respect.
type EquilibriumReaction struct {
	Reactant string
	Product1 string
	Product2 string
}

// Reaction1 is the carbon dioxide dissociation CO2 <=> CO + 0.5 O2.
var Reaction1 = EquilibriumReaction{Reactant: "CO2", Product1: "CO", Product2: "O2"}

// Reaction2 is the water dissociation H2O <=> H2 + 0.5 O2.
var Reaction2 = EquilibriumReaction{Reactant: "H2O", Product1: "H2", Product2: "O2"}

// GibbsFreeEnthalpy computes the reaction Gibbs energy dG in kJ/mol at
// temperature t for a reaction given as (reactant, product1, product2).
func GibbsFreeEnthalpy(registry *Registry, reaction EquilibriumReaction, t float64) (float64, error) {
	r, err := registry.Lookup(reaction.Reactant)
	if err != nil {
		return 0, err
	}
	p1, err := registry.Lookup(reaction.Product1)
	if err != nil {
		return 0, err
	}
	p2, err := registry.Lookup(reaction.Product2)
	if err != nil {
		return 0, err
	}
	return p1.Gibbs(t) + 0.5*p2.Gibbs(t) - r.Gibbs(t), nil
}

// EquilibriumConstant returns the pressure based equilibrium constant
//
//	Kp = exp(-dG/(RT))
//
// for the given reaction at temperature t. P is in atmospheres with a
// standard state of 1 atm.
func EquilibriumConstant(registry *Registry, reaction EquilibriumReaction, t float64) (float64, error) {
	dG, err := GibbsFreeEnthalpy(registry, reaction, t)
	if err != nil {
		return 0, err
	}
	return math.Exp(-dG * 1000.0 / (GasConstant * t)), nil
}

// KpCO2Dissociation returns Kp for CO2 <=> CO + 0.5 O2 on the default
// registry.
func KpCO2Dissociation(t float64) (float64, error) {
	kp, err := EquilibriumConstant(DefaultRegistry, Reaction1, t)
	if err != nil {
		return 0, err
	}
	return applyKp(kp, t), nil
}

// KpH2ODissociation returns Kp for H2O <=> H2 + 0.5 O2 on the default
// registry.
func KpH2ODissociation(t float64) (float64, error) {
	return EquilibriumConstant(DefaultRegistry, Reaction2, t)
}

// DissociationConstantTable returns Kp for both reactions at a temperature,
// packaged for the flame solver.
type DissociationConstantTable struct {
	Temperature     float64
	CO2Dissociation float64
	H2ODissociation float64
}

// DissociationConstants evaluates both dissociation constants once.
func DissociationConstants(t float64) (DissociationConstantTable, error) {
	k1, err := KpCO2Dissociation(t)
	if err != nil {
		return DissociationConstantTable{}, err
	}
	k2, err := KpH2ODissociation(t)
	if err != nil {
		return DissociationConstantTable{}, err
	}
	return DissociationConstantTable{Temperature: t, CO2Dissociation: k1, H2ODissociation: k2}, nil
}

// IsDissociationActive reports whether temperature t is high enough for the
// optional equilibrium refinement to change the composition noticeably.
func IsDissociationActive(t float64) bool {
	return t > 1500.0
}
