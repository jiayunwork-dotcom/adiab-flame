package thermo

import "math"

type EquilibriumReaction struct {
	Reactant string
	Product1 string
	Product2 string
}

var Reaction1 = EquilibriumReaction{Reactant: "CO2", Product1: "CO", Product2: "O2"}

var Reaction2 = EquilibriumReaction{Reactant: "H2O", Product1: "H2", Product2: "O2"}

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

func EquilibriumConstant(registry *Registry, reaction EquilibriumReaction, t float64) (float64, error) {
	dG, err := GibbsFreeEnthalpy(registry, reaction, t)
	if err != nil {
		return 0, err
	}
	return math.Exp(-dG * 1000.0 / (GasConstant * t)), nil
}

func KpCO2Dissociation(t float64) (float64, error) {
	return EquilibriumConstant(DefaultRegistry, Reaction1, t)
}

func KpH2ODissociation(t float64) (float64, error) {
	return EquilibriumConstant(DefaultRegistry, Reaction2, t)
}

type DissociationConstantTable struct {
	Temperature     float64
	CO2Dissociation float64
	H2ODissociation float64
}

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

func IsDissociationActive(t float64) bool {
	return t > 1500.0
}
