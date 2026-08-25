package flame

import "adiab-flame/internal/thermo"

type Solver struct {
	Registry *thermo.Registry
	Fuels    *thermo.FuelRegistry
}

func NewSolver() *Solver {
	return &Solver{
		Registry: thermo.DefaultRegistry,
		Fuels:    thermo.DefaultFuelRegistry,
	}
}

func (s *Solver) Solve(cfg *Config) (*Result, error) {
	cfg = cfg.WithDefaults()
	if err := cfg.Validate(s.Registry, s.Fuels); err != nil {
		return nil, err
	}
	fuel, err := s.Fuels.Lookup(cfg.Fuel)
	if err != nil {
		return nil, err
	}
	st := ComputeStoichiometry(fuel)

	reactants := ReactantStream(cfg, st)
	productBase := ProductStream(cfg, st)

	reacEnthalpy, err := reactants.Enthalpy(s.Registry, cfg.InletTemperature)
	if err != nil {
		return nil, err
	}

	res := &Result{
		Config:           *cfg,
		Fuel:             fuel.Name,
		Stoichiometry:    st,
		Reactants:        reactants,
		BaseProducts:     productBase,
		ReactantEnthalpy: reacEnthalpy,
		InletTemperature: cfg.InletTemperature,
		EquivalenceRatio: cfg.EquivalenceRatio,
		PressureAtm:      cfg.PressureAtm,
	}

	balance := func(t float64) (float64, error) {
		var products *Stream
		if cfg.Dissociation {
			products, err = solveDissociation(s.Registry, productBase, t, cfg.MaxIterations)
			if err != nil {
				return 0, err
			}
		} else {
			products = productBase.Copy()
		}
		prodEnthalpy, err := products.Enthalpy(s.Registry, t)
		if err != nil {
			return 0, err
		}
		return prodEnthalpy - reacEnthalpy, nil
	}

	tad, iterations, err := BisectTemperature(balance, cfg.InletTemperature, 3500.0, cfg.MaxIterations, cfg.TemperatureTolK)
	if err != nil {
		return nil, err
	}

	var products *Stream
	if cfg.Dissociation {
		products, err = solveDissociation(s.Registry, productBase, tad, cfg.MaxIterations)
		if err != nil {
			return nil, err
		}
	} else {
		products = productBase.Copy()
	}
	prodEnthalpy, err := products.Enthalpy(s.Registry, tad)
	if err != nil {
		return nil, err
	}
	residuals, err := AtomResiduals(s.Registry, reactants, products)
	if err != nil {
		return nil, err
	}

	res.Tad = tad
	res.Iterations = iterations
	res.ProductEnthalpy = prodEnthalpy
	res.Products = products
	res.Residuals = residuals
	res.Converged = true
	applyStoredTad(res)
	return res, nil
}

func (s *Solver) SolveSimple(fuel string, phi, tin float64) (*Result, error) {
	return s.Solve(NewConfig(fuel, phi, tin))
}
