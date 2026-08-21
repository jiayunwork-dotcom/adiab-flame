package flame

// EnthalpyBalance is the residual h_products(t) - h_reactants(Tin) in
// kJ per mole of fuel. It is monotonic increasing in t because the heat
// capacities of every product species are positive.
type EnthalpyBalance func(t float64) (float64, error)

// BisectTemperature finds the temperature where balance crosses zero on the
// bracket [tLo, tHi] using plain bisection. The bracket must satisfy
// balance(tLo) < 0 < balance(tHi). The iteration count is capped by
// maxIter and the search stops when the bracket is narrower than tolK.
func BisectTemperature(balance EnthalpyBalance, tLo, tHi float64, maxIter int, tolK float64) (float64, int, error) {
	if tLo >= tHi {
		return 0, 0, &BracketError{Low: tLo, High: tHi, Message: "empty temperature bracket"}
	}
	fLo, err := balance(tLo)
	if err != nil {
		return 0, 0, err
	}
	fHi, err := balance(tHi)
	if err != nil {
		return 0, 0, err
	}
	if fLo >= 0 {
		return 0, 0, &BracketError{Low: tLo, High: tHi, Message: "balance is already non-negative at the low edge"}
	}
	if fHi <= 0 {
		return 0, 0, &BracketError{Low: tLo, High: tHi, Message: "balance is still non-positive at the high edge"}
	}

	for i := 0; i < maxIter; i++ {
		tMid := 0.5 * (tLo + tHi)
		if (tHi-tLo)/2.0 < tolK {
			return tMid, i + 1, nil
		}
		fMid, err := balance(tMid)
		if err != nil {
			return 0, i + 1, err
		}
		if fMid == 0 {
			return tMid, i + 1, nil
		}
		if fMid < 0 {
			tLo = tMid
			fLo = fMid
		} else {
			tHi = tMid
			fHi = fMid
		}
	}
	return 0, maxIter, &NonConvergenceError{MaxIterations: maxIter, LastBracketWidth: tHi - tLo}
}

// RefineBisection restarts a bisection around a previous estimate with a
// tighter tolerance, used when the first pass only needed a coarse answer.
func RefineBisection(balance EnthalpyBalance, estimate float64, maxIter int, tolK float64) (float64, int, error) {
	width := 100.0
	lo := estimate - width
	hi := estimate + width
	fLo, err := balance(lo)
	if err != nil {
		return 0, 0, err
	}
	if fLo > 0 {
		lo = estimate - 4*width
		fLo, err = balance(lo)
		if err != nil {
			return 0, 0, err
		}
	}
	fHi, err := balance(hi)
	if err != nil {
		return 0, 0, err
	}
	if fHi < 0 {
		hi = estimate + 4*width
		fHi, err = balance(hi)
		if err != nil {
			return 0, 0, err
		}
	}
	if fLo >= 0 || fHi <= 0 {
		return 0, 0, &BracketError{
			Low:     lo,
			High:    hi,
			Message: "refinement bracket does not straddle the root",
		}
	}
	return BisectTemperature(balance, lo, hi, maxIter, tolK)
}

// residualInKelvin converts an enthalpy residual (kJ/mol) into an
// approximate temperature error using the total heat capacity of the
// products, for diagnostics.
func residualInKelvin(balance EnthalpyBalance, t float64, cpProducts float64) float64 {
	if cpProducts <= 0 {
		return 0
	}
	f, err := balance(t)
	if err != nil {
		return 0
	}
	return f * 1000.0 / cpProducts
}
